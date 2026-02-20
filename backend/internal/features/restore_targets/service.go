package restore_targets

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"databasus-backend/internal/features/databases/databases/mariadb"
	"databasus-backend/internal/features/databases/databases/mongodb"
	"databasus-backend/internal/features/databases/databases/mysql"
	"databasus-backend/internal/features/databases/databases/postgresql"
	users_models "databasus-backend/internal/features/users/models"
	workspaces_services "databasus-backend/internal/features/workspaces/services"
)

type RestoreTargetService struct {
	restoreTargetRepo *RestoreTargetRepository
	workspaceService  *workspaces_services.WorkspaceService
}

func NewRestoreTargetService(
	restoreTargetRepo *RestoreTargetRepository,
	workspaceService *workspaces_services.WorkspaceService,
) *RestoreTargetService {
	return &RestoreTargetService{
		restoreTargetRepo: restoreTargetRepo,
		workspaceService:  workspaceService,
	}
}

type CreateRestoreTargetWithDBRequest struct {
	CreateRestoreTargetRequest
	PostgresqlDatabase *postgresql.PostgresqlDatabase `json:"postgresqlDatabase"`
	MysqlDatabase      *mysql.MysqlDatabase           `json:"mysqlDatabase"`
	MariadbDatabase    *mariadb.MariadbDatabase       `json:"mariadbDatabase"`
	MongodbDatabase    *mongodb.MongodbDatabase       `json:"mongodbDatabase"`
}

func (s *RestoreTargetService) CreateRestoreTarget(
	ctx context.Context,
	user *users_models.User,
	req CreateRestoreTargetWithDBRequest,
) (*RestoreTarget, error) {
	workspace, err := s.workspaceService.GetUsersCurrentWorkspace(ctx, user)
	if err != nil {
		return nil, err
	}

	target := &RestoreTarget{
		WorkspaceID:  workspace.ID,
		Name:         req.Name,
		DatabaseType: req.DatabaseType,
		Description:  req.Description,
	}

	// Store the database reference based on type
	switch req.DatabaseType {
	case "postgresql":
		if req.PostgresqlDatabase == nil {
			return nil, errors.New("postgresql database is required")
		}
		// Hash password before storing
		if req.PostgresqlDatabase.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PostgresqlDatabase.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, err
			}
			req.PostgresqlDatabase.Password = string(hashedPassword)
		}
		// TODO: Save to databases table and get ID
		// target.PostgresDatabaseID = &postgresDB.ID

	case "mysql":
		if req.MysqlDatabase == nil {
			return nil, errors.New("mysql database is required")
		}
		if req.MysqlDatabase.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.MysqlDatabase.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, err
			}
			req.MysqlDatabase.Password = string(hashedPassword)
		}

	case "mariadb":
		if req.MariadbDatabase == nil {
			return nil, errors.New("mariadb database is required")
		}
		if req.MariadbDatabase.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.MariadbDatabase.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, err
			}
			req.MariadbDatabase.Password = string(hashedPassword)
		}

	case "mongodb":
		if req.MongodbDatabase == nil {
			return nil, errors.New("mongodb database is required")
		}
		if req.MongodbDatabase.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.MongodbDatabase.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, err
			}
			req.MongodbDatabase.Password = string(hashedPassword)
		}

	default:
		return nil, errors.New("invalid database type")
	}

	err = s.restoreTargetRepo.Create(ctx, target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

func (s *RestoreTargetService) GetRestoreTargets(
	ctx context.Context,
	user *users_models.User,
) ([]RestoreTarget, error) {
	workspace, err := s.workspaceService.GetUsersCurrentWorkspace(ctx, user)
	if err != nil {
		return nil, err
	}

	return s.restoreTargetRepo.GetByWorkspaceID(ctx, workspace.ID)
}

func (s *RestoreTargetService) GetRestoreTarget(
	ctx context.Context,
	user *users_models.User,
	id uuid.UUID,
) (*RestoreTarget, error) {
	workspace, err := s.workspaceService.GetUsersCurrentWorkspace(ctx, user)
	if err != nil {
		return nil, err
	}

	target, err := s.restoreTargetRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if target.WorkspaceID != workspace.ID {
		return nil, errors.New("restore target not found")
	}

	return target, nil
}

func (s *RestoreTargetService) UpdateRestoreTarget(
	ctx context.Context,
	user *users_models.User,
	id uuid.UUID,
	req UpdateRestoreTargetRequest,
) (*RestoreTarget, error) {
	workspace, err := s.workspaceService.GetUsersCurrentWorkspace(ctx, user)
	if err != nil {
		return nil, err
	}

	target, err := s.restoreTargetRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if target.WorkspaceID != workspace.ID {
		return nil, errors.New("restore target not found")
	}

	target.Name = req.Name
	target.Description = req.Description

	err = s.restoreTargetRepo.Update(ctx, target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

func (s *RestoreTargetService) DeleteRestoreTarget(
	ctx context.Context,
	user *users_models.User,
	id uuid.UUID,
) error {
	workspace, err := s.workspaceService.GetUsersCurrentWorkspace(ctx, user)
	if err != nil {
		return err
	}

	return s.restoreTargetRepo.Delete(ctx, id, workspace.ID)
}
