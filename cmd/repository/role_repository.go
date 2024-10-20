package repository

import (
	"gorm.io/gorm"
	"venecraft-back/cmd/entity"
)

type RoleRepository interface {
	GetRoleByName(roleName string) (*entity.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) GetRoleByName(roleName string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Where("name = ?", roleName).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetRoleNameById(id string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
