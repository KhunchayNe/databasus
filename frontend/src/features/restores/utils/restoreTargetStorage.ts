import { restoreTargetApi, type RestoreTarget as ApiRestoreTarget } from '../../../entity/restore-targets';

export type RestoreTarget = ApiRestoreTarget;

export interface StoredRestoreTarget extends RestoreTarget {
  database: any; // Full database object with credentials
}

const RESTORE_TARGETS_KEY = 'databasus_restore_targets_cache';

export const restoreTargetStorage = {
  async getTargets(): Promise<StoredRestoreTarget[]> {
    try {
      // Fetch from API
      const apiTargets = await restoreTargetApi.getRestoreTargets();

      // Load cached database credentials from localStorage
      const cachedData = localStorage.getItem(RESTORE_TARGETS_KEY);
      const cache: Record<string, any> = cachedData ? JSON.parse(cachedData) : {};

      // Merge API targets with cached credentials
      return apiTargets.map((target) => ({
        ...target,
        database: cache[target.id] || null,
      }));
    } catch (error) {
      console.error('[restoreTargetStorage] Error loading targets:', error);
      return [];
    }
  },

  async saveTarget(name: string, database: any, databaseType: string): Promise<StoredRestoreTarget> {
    try {
      console.log('[restoreTargetStorage] Saving target:', name, 'type:', databaseType);

      // Prepare request based on database type
      const request: any = {
        name,
        description: `Restore target for ${databaseType}`,
        database_type: databaseType,
      };

      // Add database-specific fields
      switch (databaseType) {
        case 'postgresql':
          request.postgresqlDatabase = database.postgresql;
          request.postgres_database_id = database.postgresql?.id;
          break;
        case 'mysql':
          request.mysqlDatabase = database.mysql;
          request.mysql_database_id = database.mysql?.id;
          break;
        case 'mariadb':
          request.mariadbDatabase = database.mariadb;
          request.mariadb_database_id = database.mariadb?.id;
          break;
        case 'mongodb':
          request.mongodbDatabase = database.mongodb;
          request.mongodb_database_id = database.mongodb?.id;
          break;
      }

      // Create restore target via API
      const newTarget = await restoreTargetApi.createRestoreTarget(request);

      // Cache the database credentials in localStorage
      const cachedData = localStorage.getItem(RESTORE_TARGETS_KEY);
      const cache: Record<string, any> = cachedData ? JSON.parse(cachedData) : {};
      cache[newTarget.id] = database;
      localStorage.setItem(RESTORE_TARGETS_KEY, JSON.stringify(cache));

      console.log('[restoreTargetStorage] Target saved successfully:', newTarget);

      // Return merged target
      return {
        ...newTarget,
        database,
      };
    } catch (error) {
      console.error('[restoreTargetStorage] Error saving target:', error);
      throw error;
    }
  },

  async deleteTarget(id: string): Promise<void> {
    try {
      console.log('[restoreTargetStorage] Deleting target:', id);

      // Delete from API
      await restoreTargetApi.deleteRestoreTarget(id);

      // Remove from localStorage cache
      const cachedData = localStorage.getItem(RESTORE_TARGETS_KEY);
      if (cachedData) {
        const cache: Record<string, any> = JSON.parse(cachedData);
        delete cache[id];
        localStorage.setItem(RESTORE_TARGETS_KEY, JSON.stringify(cache));
      }

      console.log('[restoreTargetStorage] Target deleted successfully');
    } catch (error) {
      console.error('[restoreTargetStorage] Error deleting target:', error);
      throw error;
    }
  },

  async updateTarget(id: string, name: string): Promise<void> {
    try {
      console.log('[restoreTargetStorage] Updating target:', id, 'to', name);

      // Update via API
      await restoreTargetApi.updateRestoreTarget(id, { name });

      console.log('[restoreTargetStorage] Target updated successfully');
    } catch (error) {
      console.error('[restoreTargetStorage] Error updating target:', error);
      throw error;
    }
  },
};
