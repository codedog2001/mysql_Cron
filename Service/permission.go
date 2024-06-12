package Service

import (
	"MySQL_Job/Repository"
	"context"
	"errors"
)

type PermissionService interface {
	CheckPermission(ctx context.Context, userID int64, permission string) (bool, error)
}

type permissionService struct {
	repo Repository.PermissionRepository
}

func NewPermissionService(repo Repository.PermissionRepository) PermissionService {
	return &permissionService{repo: repo}
}

// CheckPermission 鉴权的思想就是去查该用户的权限，返回权限列表，然后遍历权限，看是否有所需要的权限
func (p *permissionService) CheckPermission(ctx context.Context, userID int64, permission string) (bool, error) {
	permissions, err := p.repo.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}
	//遍历权限列表
	for _, p := range permissions {
		if p == permission {
			return true, nil
		}
	}
	return false, errors.New("permission denied")
}
