package DAO

import (
	"context"
	"gorm.io/gorm"
)

type UserDAO interface {
	GetUserPermissions(ctx context.Context, userID int64) ([]string, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewGORMUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{db: db}
}

func (dao *GORMUserDAO) GetUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	var permissions []string
	err := dao.db.WithContext(ctx).Table("permissions").
		Where("user_id = ?", userID).Pluck("permission", &permissions).Error
	return permissions, err
}
