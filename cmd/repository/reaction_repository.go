package repository

import (
	"gorm.io/gorm"
	"venecraft-back/cmd/entity"
)

type ReactionRepository interface {
	CreateReaction(reaction *entity.Reaction) error
	HasUserReacted(userID, newsID uint64, reactionType string) (bool, error)
	DeleteReaction(userID, newsID uint64, reactionType string) error
	GetReactionCountForNews(newsID uint64, reactionType string) (int64, error)
}

type reactionRepository struct {
	db *gorm.DB
}

func NewReactionRepository(db *gorm.DB) ReactionRepository {
	return &reactionRepository{db}
}

func (r *reactionRepository) CreateReaction(reaction *entity.Reaction) error {
	return r.db.Create(reaction).Error
}

func (r *reactionRepository) HasUserReacted(userID, newsID uint64, reactionType string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Reaction{}).
		Where("user_id = ? AND news_id = ? AND type = ?", userID, newsID, reactionType).
		Count(&count).Error
	return count > 0, err
}

func (r *reactionRepository) DeleteReaction(userID, newsID uint64, reactionType string) error {
	return r.db.Where("user_id = ? AND news_id = ? AND type = ?", userID, newsID, reactionType).Delete(&entity.Reaction{}).Error
}

func (r *reactionRepository) GetReactionCountForNews(newsID uint64, reactionType string) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Reaction{}).
		Where("news_id = ? AND type = ?", newsID, reactionType).
		Count(&count).Error
	return count, err
}
