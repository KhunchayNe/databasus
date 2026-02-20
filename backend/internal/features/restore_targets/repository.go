package restore_targets

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RestoreTargetRepository struct {
	db *gorm.DB
}

func NewRestoreTargetRepository(db *gorm.DB) *RestoreTargetRepository {
	return &RestoreTargetRepository{db: db}
}

func (r *RestoreTargetRepository) Create(ctx context.Context, target *RestoreTarget) error {
	target.ID = uuid.New()
	target.CreatedAt = time.Now()
	target.UpdatedAt = time.Now()

	return r.db.WithContext(ctx).Create(target).Error
}

func (r *RestoreTargetRepository) GetByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]RestoreTarget, error) {
	var targets []RestoreTarget
	err := r.db.WithContext(ctx).
		Where("workspace_id = ?", workspaceID).
		Order("name ASC").
		Find(&targets).Error

	if err != nil {
		return nil, err
	}

	return targets, nil
}

func (r *RestoreTargetRepository) GetByID(ctx context.Context, id uuid.UUID) (*RestoreTarget, error) {
	var target RestoreTarget
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&target).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("restore target not found")
		}
		return nil, err
	}

	return &target, nil
}

func (r *RestoreTargetRepository) Update(ctx context.Context, target *RestoreTarget) error {
	target.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).
		Model(&RestoreTarget{}).
		Where("id = ? AND workspace_id = ?", target.ID, target.WorkspaceID).
		Updates(target)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("restore target not found")
	}

	return nil
}

func (r *RestoreTargetRepository) Delete(ctx context.Context, id uuid.UUID, workspaceID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND workspace_id = ?", id, workspaceID).
		Delete(&RestoreTarget{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("restore target not found")
	}

	return nil
}
