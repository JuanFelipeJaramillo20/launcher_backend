package repository

import (
	"gorm.io/gorm"
	"venecraft-back/cmd/entity"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	GetAllUsers() ([]entity.User, error)
	GetUserByID(id uint64) (*entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteUser(id uint64) error
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByNickname(nickname string) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetAllUsers() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Preload("Roles").Find(&users).Error
	return users, err
}

func (r *userRepository) GetUserByID(id uint64) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Roles").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) DeleteUser(id uint64) error {
	return r.db.Delete(&entity.User{}, id).Error
}

func (r *userRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Roles").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByNickname(nickname string) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Roles").Where("nickname = ?", nickname).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
