package service

import (
	"errors"
	"fmt"
	"time"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/repository"
)

type ReactionService interface {
	ToggleReaction(userID, newsID uint64, reactionType string) (bool, error)
}

type reactionService struct {
	reactionRepo repository.ReactionRepository
	logRepo      repository.LogRepository
}

func NewReactionService(reactionRepo repository.ReactionRepository, logRepo repository.LogRepository) ReactionService {
	return &reactionService{reactionRepo, logRepo}
}

func (s *reactionService) ToggleReaction(userID, newsID uint64, reactionType string) (bool, error) {
	if reactionType != "like" && reactionType != "dislike" {
		return false, errors.New("invalid reaction type")
	}

	// Check if the user already reacted with the specified type
	reacted, err := s.reactionRepo.HasUserReacted(userID, newsID, reactionType)
	if err != nil {
		return false, err
	}

	if reacted {
		// If the user has already reacted, remove the reaction
		err := s.reactionRepo.DeleteReaction(userID, newsID, reactionType)
		if err != nil {
			return false, err
		}

		logEntry := entity.Log{
			UserID:      userID,
			Action:      fmt.Sprintf("un%s", reactionType),
			Description: fmt.Sprintf("User with id: %d un%sd the news post with id: %d", userID, reactionType, newsID),
			Timestamp:   time.Now(),
		}
		if logErr := s.logRepo.CreateLog(&logEntry); logErr != nil {
			return false, logErr
		}

		return false, nil
	}

	// Add the new reaction
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
