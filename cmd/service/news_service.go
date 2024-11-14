package service

import (
	"errors"
	"fmt"
	"time"
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
	ToggleLikeNews(userID, newsID uint64) (bool, error)
}

type newsService struct {
	newsRepo repository.NewsRepository
	likeRepo repository.LikeRepository
	logRepo  repository.LogRepository
}

func NewNewsService(newsRepo repository.NewsRepository, likeRepo repository.LikeRepository, logRepo repository.LogRepository) NewsService {
	return &newsService{newsRepo, likeRepo, logRepo}
}

func (s *newsService) CreateNews(news *entity.News) error {
	if news.Title == "" || news.Content == "" {
		return errors.New("the title and content must be provided")
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

func (s *newsService) ToggleLikeNews(userID, newsID uint64) (bool, error) {
	liked, err := s.likeRepo.HasUserLikedNews(userID, newsID)
	if err != nil {
		return false, err
	}

	if liked {
		err := s.likeRepo.DeleteLike(userID, newsID)
		if err != nil {
			return false, err
		}

		logEntry := entity.Log{
			UserID:      userID,
			Action:      "unlike",
			Description: fmt.Sprintf("User with id: %d unliked the news post with id: %d", userID, newsID),
			Timestamp:   time.Now(),
		}
		if logErr := s.logRepo.CreateLog(&logEntry); logErr != nil {
			return false, logErr
		}

		return false, nil
	}

	like := entity.Like{
		UserID:    userID,
		NewsID:    newsID,
		Timestamp: time.Now(),
	}
	if err := s.likeRepo.CreateLike(&like); err != nil {
		return false, err
	}

	logEntry := entity.Log{
		UserID:      userID,
		Action:      "like",
		Description: fmt.Sprintf("User with id: %d liked the news post with id: %d", userID, newsID),
		Timestamp:   time.Now(),
	}
	if logErr := s.logRepo.CreateLog(&logEntry); logErr != nil {
		return false, logErr
	}

	return true, nil
}
