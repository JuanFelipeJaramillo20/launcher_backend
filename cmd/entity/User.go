package entity

import "time"

// swagger:model User
type User struct {
	// User ID
	// required: true
	ID uint64 `json:"id" gorm:"primaryKey;autoIncrement"`

	// Full name of the user
	// required: true
	FullName string `json:"full_name" gorm:"type:varchar(255)"`

	// Email address
	// required: true
	Email string `json:"email" gorm:"type:varchar(255);unique"`

	// Nickname of the user
	// required: true
	Nickname string `json:"nickname" gorm:"type:varchar(100);unique"`

	// Password (not included in responses)
	// required: true
	Password string `json:"-" gorm:"type:varchar(255)"`

	// Token for password recovery
	RecoverPasswordToken string `json:"-" gorm:"type:varchar(255)"`

	// Expiration time of password recovery token
	RecoverPasswordTokenExpires time.Time `gorm:"default:null"`

	// List of user roles
	Roles []*Role `gorm:"many2many:user_roles;-" json:"roles,omitempty"`

	// Active status of the user
	IsActive bool `json:"is_active" gorm:"default:true"`
}
