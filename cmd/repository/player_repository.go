package repository

import (
	"gorm.io/gorm"
	"venecraft-back/cmd/entity"
)

type PlayerRepository interface {
	CreatePlayer(player *entity.Player) error
	UpdatePlayer(player *entity.Player) error
	DeletePlayer(playerID uint64) error
	GetPlayerByID(playerID uint64) (*entity.Player, error)
}

type playerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &playerRepository{db}
}

func (r *playerRepository) CreatePlayer(player *entity.Player) error {
	return r.db.Create(player).Error
}

func (r *playerRepository) UpdatePlayer(player *entity.Player) error {
	return r.db.Save(player).Error
}

func (r *playerRepository) DeletePlayer(playerID uint64) error {
	return r.db.Delete(&entity.Player{}, playerID).Error
}

func (r *playerRepository) GetPlayerByID(playerID uint64) (*entity.Player, error) {
	var player entity.Player
	err := r.db.First(&player, playerID).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}
