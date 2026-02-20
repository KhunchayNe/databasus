import type { PostgresqlDatabase } from '../../databases/model/postgresql/PostgresqlDatabase';
import type { MysqlDatabase } from '../../databases/model/mysql/MysqlDatabase';
import type { MariadbDatabase } from '../../databases/model/mariadb/MariadbDatabase';
import type { MongodbDatabase } from '../../databases/model/mongodb/MongodbDatabase';

export interface RestoreTarget {
  id: string;
  workspace_id: string;
  name: string;
  database_type: 'postgresql' | 'mysql' | 'mariadb' | 'mongodb';
  description?: string;

  // PostgreSQL fields
  postgres_database_id?: string;
  postgresqlDatabase?: PostgresqlDatabase;

  // MySQL fields
  mysql_database_id?: string;
  mysqlDatabase?: MysqlDatabase;

  // MariaDB fields
  mariadb_database_id?: string;
  mariadbDatabase?: MariadbDatabase;

  // MongoDB fields
  mongodb_database_id?: string;
  mongodbDatabase?: MongodbDatabase;

  created_at: string;
  updated_at: string;
}
