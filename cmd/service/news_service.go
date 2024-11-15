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
	ToggleReactionNews(userID, newsID uint64, reactionType string) (bool, error)
}

type newsService struct {
	newsRepo     repository.NewsRepository
	reactionRepo repository.ReactionRepository
	logRepo      repository.LogRepository
}

func NewNewsService(newsRepo repository.NewsRepository, reactionRepo repository.ReactionRepository, logRepo repository.LogRepository) NewsService {
	return &newsService{newsRepo, reactionRepo, logRepo}
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
		likeCount, err := s.reactionRepo.GetReactionCountForNews(news.ID, "like")
		if err != nil {
			return nil, err
		}

		dislikeCount, err := s.reactionRepo.GetReactionCountForNews(news.ID, "dislike")
		if err != nil {
			return nil, err
		}

		userLiked, err := s.reactionRepo.HasUserReacted(userID, news.ID, "like")
		if err != nil {
			return nil, err
		}

		userDisliked, err := s.reactionRepo.HasUserReacted(userID, news.ID, "dislike")
		if err != nil {
			return nil, err
		}

		newsData := map[string]interface{}{
			"id":            news.ID,
			"title":         news.Title,
			"content":       news.Content,
			"created_by":    news.CreatedBy,
			"created_at":    news.CreatedAt,
			"like_count":    likeCount,
			"dislike_count": dislikeCount,
			"user_liked":    userLiked,
			"user_disliked": userDisliked,
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
		likeCount, err := s.reactionRepo.GetReactionCountForNews(news.ID, "like")
		if err != nil {
			return nil, err
		}

		dislikeCount, err := s.reactionRepo.GetReactionCountForNews(news.ID, "dislike")
		if err != nil {
			return nil, err
		}

		userLiked, err := s.reactionRepo.HasUserReacted(userID, news.ID, "like")
		if err != nil {
			return nil, err
		}

		userDisliked, err := s.reactionRepo.HasUserReacted(userID, news.ID, "dislike")
		if err != nil {
			return nil, err
		}

		newsData := map[string]interface{}{
			"id":            news.ID,
			"title":         news.Title,
			"content":       news.Content,
			"created_by":    news.CreatedBy,
			"created_at":    news.CreatedAt,
			"like_count":    likeCount,
			"dislike_count": dislikeCount,
			"user_liked":    userLiked,
			"user_disliked": userDisliked,
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

	likeCount, err := s.reactionRepo.GetReactionCountForNews(news.ID, "like")
	if err != nil {
		return nil, err
	}

	dislikeCount, err := s.reactionRepo.GetReactionCountForNews(news.ID, "dislike")
	if err != nil {
		return nil, err
	}

	userLiked, err := s.reactionRepo.HasUserReacted(userID, news.ID, "like")
	if err != nil {
		return nil, err
	}

	userDisliked, err := s.reactionRepo.HasUserReacted(userID, news.ID, "dislike")
	if err != nil {
		return nil, err
	}

	newsData := map[string]interface{}{
		"id":            news.ID,
		"title":         news.Title,
		"content":       news.Content,
		"created_by":    news.CreatedBy,
		"created_at":    news.CreatedAt,
		"like_count":    likeCount,
		"dislike_count": dislikeCount,
		"user_liked":    userLiked,
		"user_disliked": userDisliked,
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

func (s *newsService) ToggleReactionNews(userID, newsID uint64, reactionType string) (bool, error) {
	reacted, err := s.reactionRepo.HasUserReacted(userID, newsID, reactionType)
	if err != nil {
		return false, err
	}

	if reacted {
		err := s.reactionRepo.DeleteReaction(userID, newsID, reactionType)
		if err != nil {
			return false, err
		}

		logEntry := entity.Log{
			UserID:      userID,
			Action:      "un" + reactionType,
			Description: fmt.Sprintf("User with id: %d un%sd the news post with id: %d", userID, reactionType, newsID),
			Timestamp:   time.Now(),
		}
		if logErr := s.logRepo.CreateLog(&logEntry); logErr != nil {
			return false, logErr
		}

		return false, nil
	}

	reaction := entity.Reaction{
		UserID:    userID,
		NewsID:    newsID,
		Type:      reactionType,
		Timestamp: time.Now(),
	}
	if err := s.reactionRepo.CreateReaction(&reaction); err != nil {
		return false, err
	}

	logEntry := entity.Log{
		UserID:      userID,
		Action:      reactionType,
		Description: fmt.Sprintf("User with id: %d %sd the news post with id: %d", userID, reactionType, newsID),
		Timestamp:   time.Now(),
	}
	if logErr := s.logRepo.CreateLog(&logEntry); logErr != nil {
		return false, logErr
	}

	return true, nil
}
