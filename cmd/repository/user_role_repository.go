package repository

import (
	"gorm.io/gorm"
	"venecraft-back/cmd/entity"
)

type UserRoleRepository interface {
	AssignRole(userRole *entity.UserRole) error
}

type userRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) UserRoleRepository {
	return &userRoleRepository{db}
}

func (r *userRoleRepository) AssignRole(userRole *entity.UserRole) error {
	return r.db.Create(userRole).Error
}
