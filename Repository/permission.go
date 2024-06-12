package Repository

import (
	"MySQL_Job/Repository/DAO"
	"context"
)

type PermissionRepository interface {
	GetUserPermissions(ctx context.Context, userID int64) ([]string, error)
}

type permissionRepository struct {
	dao DAO.UserDAO
}

func NewPermissionRepository(dao DAO.UserDAO) PermissionRepository {
	return &permissionRepository{dao: dao}
}

func (repo *permissionRepository) GetUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	return repo.dao.GetUserPermissions(ctx, userID)
}
