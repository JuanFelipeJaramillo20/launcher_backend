package service

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/repository"
)

type NewsService interface {
	CreateNews(news *entity.News) error
	GetAllNews() ([]entity.News, error)
	GetNewsByID(id uint64) (*entity.User, error)
	UpdateNews(news *entity.News) error
	DeleteNews(id uint64) error
}