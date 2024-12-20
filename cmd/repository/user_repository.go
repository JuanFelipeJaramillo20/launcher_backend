package repository

import (
	"fmt"
	"gorm.io/gorm"
	"venecraft-back/cmd/entity"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	GetAllUsers() ([]entity.User, error)
	GetUserByID(id uint64) (*entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteUser(id uint64) error
	GetUserByEmail(email string, preloadRoles bool) (*entity.User, error)
	GetUserByNickname(nickname string) (*entity.User, error)
	GetUsersByRole(role string) ([]entity.User, error)
	GetUserByResetToken(resetToken string) (*entity.User, error)
	HasRole(id uint64, role string) bool
	CountActiveUsers() (int, error)
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
	return r.db.Model(&entity.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"full_name":                      user.FullName,
			"email":                          user.Email,
			"nickname":                       user.Nickname,
			"password":                       user.Password,
			"recover_password_token":         user.RecoverPasswordToken,
			"recover_password_token_expires": user.RecoverPasswordTokenExpires,
			"is_active":                      user.IsActive,
		}).Error
}

func (r *userRepository) DeleteUser(id uint64) error {
	return r.db.Delete(&entity.User{}, id).Error
}

func (r *userRepository) GetUserByEmail(email string, preloadRoles bool) (*entity.User, error) {
	var user entity.User
	query := r.db.Where("email = ?", email)

	// Conditionally preload roles
	if preloadRoles {
		query = query.Preload("Roles")
	}

	err := query.First(&user).Error
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

// GetUsersByRole retrieves all users assigned to a specific role.
func (r *userRepository) GetUsersByRole(role string) ([]entity.User, error) {
	var users []entity.User
	err := r.db.Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("roles.name = ?", role).
		Preload("Roles").
		Find(&users).Error
	return users, err
}

func (r *userRepository) GetUserByResetToken(resetToken string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("recover_password_token = ?", resetToken).First(&user).Error
	return &user, err
}

func (r *userRepository) HasRole(userID uint64, roleName string) bool {
	var count int64
	err := r.db.Table("user_roles").
		Select("count(*)").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.name = ?", userID, roleName).
		Count(&count).Error

	if err != nil {
		fmt.Printf("Error checking role for user %d: %v\n", userID, err)
		return false
	}

	return count > 0
}

func (r *userRepository) CountActiveUsers() (int, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Where("is_active = ?", true).Count(&count).Error
	return int(count), err
}
