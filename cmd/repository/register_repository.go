package repository

import (
	"gorm.io/gorm"
	"venecraft-back/cmd/entity"
)

type RegisterRepository interface {
	CreateRegister(register *entity.Register) error
	GetRegisterByID(id uint64) (*entity.Register, error)
	DeleteRegister(id uint64) error
	UpdateRegister(register *entity.Register) error
}

type registerRepository struct {
	db *gorm.DB
}

func NewRegisterRepository(db *gorm.DB) RegisterRepository {
	return &registerRepository{db}
}

func (r *registerRepository) CreateRegister(register *entity.Register) error {
	return r.db.Create(register).Error
}

func (r *registerRepository) GetRegisterByID(id uint64) (*entity.Register, error) {
	var reg entity.Register
	err := r.db.First(&reg, id).Error
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *registerRepository) DeleteRegister(id uint64) error {
	return r.db.Delete(&entity.Register{}, id).Error
}

func (r *registerRepository) UpdateRegister(register *entity.Register) error {
	return r.db.Save(register).Error
}
