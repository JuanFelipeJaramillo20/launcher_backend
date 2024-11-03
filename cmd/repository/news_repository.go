package repository

import (
	"errors"

	"venecraft-back/cmd/entity"

	"gorm.io/gorm"
)

type NewsRepository interface {
	CreateNews(news *entity.News) error
	GetNewsByID(id uint64) (*entity.News, error)
	GetAllNews() ([]entity.News, error)
	GetLatestNews() ([]entity.News, error)
	UpdateNews(news *entity.News) error
	DeleteNews(id uint64) error
}

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) NewsRepository {
	return &newsRepository{db}
}

func (r *newsRepository) CreateNews(news *entity.News) error {
	if err := r.db.Create(news).Error; err != nil {
		return err
	}
	return nil
}

func (r *newsRepository) GetNewsByID(id uint64) (*entity.News, error) {
	var news entity.News
	if err := r.db.First(&news, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &news, nil
}

func (r *newsRepository) GetAllNews() ([]entity.News, error) {
	var news []entity.News
	if err := r.db.Find(&news).Error; err != nil {
		return nil, err
	}
	return news, nil
}

func (r *newsRepository) GetLatestNews() ([]entity.News, error) {
	var news []entity.News
	if err := r.db.Order("created_at DESC").Find(&news).Error; err != nil {
		return nil, err
	}
	return news, nil
}

func (r *newsRepository) UpdateNews(news *entity.News) error {
	return r.db.Save(news).Error
}

func (r *newsRepository) DeleteNews(id uint64) error {
	if err := r.db.Delete(&entity.News{}, id).Error; err != nil {
		return err
	}
	return nil
}
