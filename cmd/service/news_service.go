package service

import (
	"errors"
	"fmt"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/repository"
)

type NewsService interface {
	CreateNews(news *entity.News) error
	GetAllNews() ([]entity.News, error)
	GetLatestNews() ([]entity.News, error)
	GetNewsByID(id uint64) (*entity.News, error)
	UpdateNews(news *entity.News) error
	DeleteNews(id uint64) error
}

type newsService struct {
	newsRepo repository.NewsRepository
}

func NewNewsService(newsRepo repository.NewsRepository) NewsService {
	return &newsService{newsRepo}
}

func (s *newsService) CreateNews(news *entity.News) error {
	if news.Title == "" || news.Content == "" {
		return errors.New("the tittle and content must be provided")
	}
	return s.newsRepo.CreateNews(news)
}

func (s *newsService) GetAllNews() ([]entity.News, error) {
	return s.newsRepo.GetAllNews()
}

func (s *newsService) GetLatestNews() ([]entity.News, error) {
	return s.newsRepo.GetLatestNews()
}

func (s *newsService) GetNewsByID(id uint64) (*entity.News, error) {
	news, err := s.newsRepo.GetNewsByID(id)
	if err != nil {
		return nil, err
	}
	if news == nil {
		return nil, errors.New("news not found")
	}
	return news, nil
}

func (s *newsService) UpdateNews(news *entity.News) error {
	fmt.Print(news)
	fmt.Print(news.ID)
	return s.newsRepo.UpdateNews(news)
}

func (s *newsService) DeleteNews(id uint64) error {
	return s.newsRepo.DeleteNews(id)
}
