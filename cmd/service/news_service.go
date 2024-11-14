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
	GetAllNews(userID uint64) ([]map[string]interface{}, error)
	GetLatestNews(userID uint64) ([]map[string]interface{}, error)
	GetNewsByID(userID uint64, id uint64) (map[string]interface{}, error)
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

func (s *newsService) GetAllNews(userID uint64) ([]map[string]interface{}, error) {
	newsList, err := s.newsRepo.GetAllNews()
	if err != nil {
		return nil, err
	}

	var response []map[string]interface{}
	for _, news := range newsList {
		likeCount, err := s.newsRepo.GetLikeCountForNews(news.ID)
		if err != nil {
			return nil, err
		}

		userLiked, err := s.likeRepo.HasUserLikedNews(userID, news.ID)
		if err != nil {
			return nil, err
		}

		newsData := map[string]interface{}{
			"id":         news.ID,
			"title":      news.Title,
			"content":    news.Content,
			"created_by": news.CreatedBy,
			"created_at": news.CreatedAt,
			"like_count": likeCount,
			"user_liked": userLiked,
		}

		response = append(response, newsData)
	}

	return response, nil
}

func (s *newsService) GetLatestNews(userID uint64) ([]map[string]interface{}, error) {
	newsList, err := s.newsRepo.GetLatestNews()
	if err != nil {
		return nil, err
	}

	var response []map[string]interface{}
	for _, news := range newsList {
		likeCount, err := s.newsRepo.GetLikeCountForNews(news.ID)
		if err != nil {
			return nil, err
		}

		userLiked, err := s.likeRepo.HasUserLikedNews(userID, news.ID)
		if err != nil {
			return nil, err
		}

		newsData := map[string]interface{}{
			"id":         news.ID,
			"title":      news.Title,
			"content":    news.Content,
			"created_by": news.CreatedBy,
			"created_at": news.CreatedAt,
			"like_count": likeCount,
			"user_liked": userLiked,
		}

		response = append(response, newsData)
	}

	return response, nil
}

func (s *newsService) GetNewsByID(userID, newsID uint64) (map[string]interface{}, error) {
	news, err := s.newsRepo.GetNewsByID(newsID)
	if err != nil {
		return nil, err
	}
	if news == nil {
		return nil, errors.New("news not found")
	}

	likeCount, err := s.newsRepo.GetLikeCountForNews(news.ID)
	if err != nil {
		return nil, err
	}

	userLiked, err := s.likeRepo.HasUserLikedNews(userID, news.ID)
	if err != nil {
		return nil, err
	}

	newsData := map[string]interface{}{
		"id":         news.ID,
		"title":      news.Title,
		"content":    news.Content,
		"created_by": news.CreatedBy,
		"created_at": news.CreatedAt,
		"like_count": likeCount,
		"user_liked": userLiked,
	}

	return newsData, nil
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
