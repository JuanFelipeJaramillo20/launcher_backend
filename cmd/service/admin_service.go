package service

import (
	"errors"
	"time"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/repository"
)

type AdminService interface {
	CreatePlayer(player *entity.Player) error
	UpdatePlayer(player *entity.Player) error
	DeletePlayer(playerID uint64) error
	GetPlayerByID(playerID uint64) (*entity.Player, error)

	CreateModerator(mod *entity.User) error
	UpdateModerator(mod *entity.User) error
	DeleteModerator(modID uint64) error
	GetModeratorByID(modID uint64) (*entity.User, error)

	BanPlayer(playerID uint64, reason string, duration time.Duration) error
}

type adminService struct {
	userRepo    repository.UserRepository
	playerRepo  repository.PlayerRepository
	banRepo     repository.BanRepository
	roleRepo    repository.RoleRepository
	userService UserService
}

func NewAdminService(userRepo repository.UserRepository, playerRepo repository.PlayerRepository, banRepo repository.BanRepository, roleRepo repository.RoleRepository, userService UserService) AdminService {
	return &adminService{
		userRepo:    userRepo,
		playerRepo:  playerRepo,
		banRepo:     banRepo,
		roleRepo:    roleRepo,
		userService: userService,
	}
}

func (s *adminService) CreatePlayer(player *entity.User) error {
	return s.userService.CreateUser(player, "PLAYER")
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
	// Verify the user has the moderator role
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
