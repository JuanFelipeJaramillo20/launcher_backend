package repository

import (
	"gorm.io/gorm"
	"venecraft-back/cmd/entity"
)

type BanRepository interface {
	CreateBan(ban *entity.Ban) error
	GetBanByPlayerID(playerID uint64) (*entity.Ban, error)
	UpdateBan(ban *entity.Ban) error
	DeleteBan(banID uint64) error
}

type banRepository struct {
	db *gorm.DB
}

func NewBanRepository(db *gorm.DB) BanRepository {
	return &banRepository{db}
}

func (r *banRepository) CreateBan(ban *entity.Ban) error {
	return r.db.Create(ban).Error
}

func (r *banRepository) GetBanByPlayerID(playerID uint64) (*entity.Ban, error) {
	var ban entity.Ban
	err := r.db.Where("player_id = ?", playerID).First(&ban).Error
	if err != nil {
		return nil, err
	}
	return &ban, nil
}

func (r *banRepository) UpdateBan(ban *entity.Ban) error {
	return r.db.Save(ban).Error
}

func (r *banRepository) DeleteBan(banID uint64) error {
	return r.db.Delete(&entity.Ban{}, banID).Error
}
