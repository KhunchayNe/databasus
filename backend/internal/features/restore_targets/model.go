package restore_targets

import (
	"time"

	"github.com/google/uuid"
)

type RestoreTarget struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	WorkspaceID uuid.UUID              `json:"workspace_id" db:"workspace_id"`
	Name        string                 `json:"name" db:"name"`
	DatabaseType string                `json:"database_type" db:"database_type"`
	Description string                 `json:"description" db:"description"`

	// PostgreSQL fields
	PostgresDatabaseID *uuid.UUID       `json:"postgres_database_id,omitempty" db:"postgres_database_id"`

	// MySQL fields
	MysqlDatabaseID *uuid.UUID          `json:"mysql_database_id,omitempty" db:"mysql_database_id"`

	// MariaDB fields
	MariadbDatabaseID *uuid.UUID         `json:"mariadb_database_id,omitempty" db:"mariadb_database_id"`

	// MongoDB fields
	MongodbDatabaseID *uuid.UUID         `json:"mongodb_database_id,omitempty" db:"mongodb_database_id"`

	CreatedAt time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at" db:"updated_at"`
}

type CreateRestoreTargetRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	DatabaseType string `json:"database_type" binding:"required,oneof=postgresql mysql mariadb mongodb"`
}

type UpdateRestoreTargetRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
