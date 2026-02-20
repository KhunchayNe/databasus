export interface CreateRestoreTargetRequest {
  name: string;
  description?: string;
  database_type: 'postgresql' | 'mysql' | 'mariadb' | 'mongodb';
}

export interface CreateRestoreTargetWithDBRequest extends CreateRestoreTargetRequest {
  postgresqlDatabase?: Record<string, unknown>;
  mysqlDatabase?: Record<string, unknown>;
  mariadbDatabase?: Record<string, unknown>;
  mongodbDatabase?: Record<string, unknown>;
}

export interface UpdateRestoreTargetRequest {
  name?: string;
  description?: string;
}
