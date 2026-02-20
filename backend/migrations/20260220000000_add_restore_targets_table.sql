-- +goose Up
START TRANSACTION;

CREATE TABLE restore_targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    database_type VARCHAR(50) NOT NULL, -- 'postgresql', 'mysql', 'mariadb', 'mongodb'
    description TEXT,

    -- PostgreSQL fields
    pg_host VARCHAR(255),
    pg_port INTEGER,
    pg_database VARCHAR(255),
    pg_username VARCHAR(255),
    pg_password TEXT,
    pg_version VARCHAR(10),
    pg_ssl_mode VARCHAR(20),
    pg_exclude_extensions BOOLEAN DEFAULT false,
    pg_include_schemas TEXT,

    -- MySQL fields
    mysql_host VARCHAR(255),
    mysql_port INTEGER,
    mysql_database VARCHAR(255),
    mysql_username VARCHAR(255),
    mysql_password TEXT,
    mysql_version VARCHAR(10),
    mysql_ssl_mode BOOLEAN DEFAULT false,
    mysql_is_exclude_events BOOLEAN DEFAULT false,
    mysql_privileges TEXT,

    -- MariaDB fields
    mariadb_host VARCHAR(255),
    mariadb_port INTEGER,
    mariadb_database VARCHAR(255),
    mariadb_username VARCHAR(255),
    mariadb_password TEXT,
    mariadb_version VARCHAR(10),
    mariadb_client_version VARCHAR(10),
    mariadb_ssl_mode BOOLEAN DEFAULT false,
    mariadb_is_exclude_events BOOLEAN DEFAULT false,
    mariadb_privileges TEXT,

    -- MongoDB fields
    mongodb_host VARCHAR(255),
    mongodb_port INTEGER,
    mongodb_database VARCHAR(255),
    mongodb_username VARCHAR(255),
    mongodb_password TEXT,
    mongodb_version VARCHAR(10),
    mongodb_auth_database VARCHAR(255),
    mongodb_ssl_mode BOOLEAN DEFAULT false,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster queries
CREATE INDEX idx_restore_targets_workspace_id ON restore_targets(workspace_id);
CREATE INDEX idx_restore_targets_database_type ON restore_targets(database_type);
CREATE INDEX idx_restore_targets_name ON restore_targets(name);

COMMIT;

-- +goose Down
START TRANSACTION;

DROP INDEX IF EXISTS idx_restore_targets_name;
DROP INDEX IF EXISTS idx_restore_targets_database_type;
DROP INDEX IF EXISTS idx_restore_targets_workspace_id;
DROP TABLE IF EXISTS restore_targets;

COMMIT;
