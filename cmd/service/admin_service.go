package service

import (
	"errors"
	"time"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/repository"
)

type AdminService interface {
	CreateUser(user *entity.User, role string) error
	UpdateUser(user *entity.User) error
	DeleteUser(userID uint64) error
	GetUserByID(userID uint64) (*entity.User, error)

	CreatePlayer(player *entity.Player) error
	UpdatePlayer(player *entity.Player) error
	DeletePlayer(playerID uint64) error
	GetPlayerByID(playerID uint64) (*entity.Player, error)

	CreateModerator(mod *entity.User) error
	UpdateModerator(mod *entity.User) error
	DeleteModerator(modID uint64) error
	GetModeratorByID(modID uint64) (*entity.User, error)

	BanPlayer(playerID uint64, reason string, duration time.Duration) error

	CreateNews(news *entity.News) error
	GetAllNews() ([]entity.News, error)
	GetLatestNews() ([]entity.News, error)
	GetNewsByID(id uint64) (*entity.News, error)
	UpdateNews(news *entity.News) error
	DeleteNews(id uint64) error
}

type adminService struct {
	userRepo    repository.UserRepository
	newsRepo    repository.NewsRepository
	playerRepo  repository.PlayerRepository
	banRepo     repository.BanRepository
	roleRepo    repository.RoleRepository
	userService UserService
	newsService NewsService
}

func NewAdminService(userRepo repository.UserRepository, newsRepo repository.NewsRepository, playerRepo repository.PlayerRepository, banRepo repository.BanRepository, roleRepo repository.RoleRepository, userService UserService, newsService NewsService) AdminService {
	return &adminService{
		userRepo:    userRepo,
		newsRepo:    newsRepo,
		playerRepo:  playerRepo,
		banRepo:     banRepo,
		roleRepo:    roleRepo,
		userService: userService,
		newsService: newsService,
	}
}

func (s *adminService) CreateUser(user *entity.User, role string) error {
	return s.userService.CreateUser(user, role)
}

func (s *adminService) UpdateUser(user *entity.User) error {
	return s.userRepo.UpdateUser(user)
}

func (s *adminService) DeleteUser(userID uint64) error {
	return s.userRepo.DeleteUser(userID)
}

func (s *adminService) GetUserByID(userID uint64) (*entity.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *adminService) CreatePlayer(player *entity.Player) error {
	return s.playerRepo.CreatePlayer(player)
}

func (s *adminService) UpdatePlayer(player *entity.Player) error {
	return s.playerRepo.UpdatePlayer(player)
}

func (s *adminService) DeletePlayer(playerID uint64) error {
	return s.playerRepo.DeletePlayer(playerID)
}

func (s *adminService) GetPlayerByID(playerID uint64) (*entity.Player, error) {
	player, err := s.playerRepo.GetPlayerByID(playerID)
	if err != nil {
		return nil, errors.New("player not found")
	}
	return player, nil
}

func (s *adminService) CreateModerator(mod *entity.User) error {
	return s.userService.CreateUser(mod, "MODERATOR")
}

func (s *adminService) UpdateModerator(mod *entity.User) error {
	if !s.userRepo.HasRole(mod.ID, "MODERATOR") {
		return errors.New("user is not a moderator")
	}
	return s.userRepo.UpdateUser(mod)
}

func (s *adminService) DeleteModerator(modID uint64) error {
	user, err := s.userRepo.GetUserByID(modID)
	if err != nil || !s.userRepo.HasRole(modID, "MODERATOR") {
		return errors.New("moderator not found")
	}
	return s.userRepo.DeleteUser(user.ID)
}

func (s *adminService) GetModeratorByID(modID uint64) (*entity.User, error) {
	user, err := s.userRepo.GetUserByID(modID)
	if err != nil || !s.userRepo.HasRole(modID, "MODERATOR") {
		return nil, errors.New("moderator not found")
	}
	return user, nil
}

func (s *adminService) BanPlayer(playerID uint64, reason string, duration time.Duration) error {
	player, err := s.playerRepo.GetPlayerByID(playerID)
	if err != nil {
		return errors.New("player not found")
	}

	ban := &entity.Ban{
		PlayerID: player.ID,
		Reason:   reason,
		BanDate:  time.Now(),
		Duration: duration,
	}

	err = s.banRepo.CreateBan(ban)
	if err != nil {
		return errors.New("failed to create ban")
	}

	return nil
}

func (s *adminService) CreateNews(news *entity.News) error {
	return s.newsService.CreateNews(news)
}

func (s *adminService) GetAllNews() ([]entity.News, error) {
	return s.newsRepo.GetAllNews()
}

func (s *adminService) GetLatestNews() ([]entity.News, error) {
	return s.newsRepo.GetLatestNews()
}

func (s *adminService) GetNewsByID(newsID uint64) (*entity.News, error) {
	news, err := s.newsRepo.GetNewsByID(newsID)
	if err != nil {
		return nil, errors.New("news not found")
	}
	return news, nil
}

func (s *adminService) UpdateNews(news *entity.News) error {
	return s.newsRepo.UpdateNews(news)
}

func (s *adminService) DeleteNews(newsID uint64) error {
	return s.newsRepo.DeleteNews(newsID)
}
