package repository

import (
	"gorm.io/gorm"
	"venecraft-back/cmd/entity"
)

type LikeRepository interface {
	CreateLike(like *entity.Like) error
	HasUserLikedNews(userID, newsID uint64) (bool, error)
	DeleteLike(userID, newsID uint64) error
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db}
}

func (r *likeRepository) CreateLike(like *entity.Like) error {
	return r.db.Create(like).Error
}

func (r *likeRepository) HasUserLikedNews(userID, newsID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Like{}).
		Where("user_id = ? AND news_id = ?", userID, newsID).
		Count(&count).Error
	return count > 0, err
}

func (r *likeRepository) DeleteLike(userID, newsID uint64) error {
	return r.db.Where("user_id = ? AND news_id = ?", userID, newsID).Delete(&entity.Like{}).Error
}
