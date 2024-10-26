package entity

// swagger:model Register
type Register struct {
	// Register ID
	// required: true
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// Full name for registration
	// required: true
	FullName string `gorm:"type:varchar(255)"`

	// Email address for registration
	// required: true
	Email string `gorm:"type:varchar(255);unique"`

	// Nickname chosen for registration
	// required: true
	Nickname string `gorm:"type:varchar(100);unique"`

	// Password for registration
	// required: true
	Password string `gorm:"type:varchar(255)"`

	// Status of registration approval
	Accepted bool `gorm:"default:false"`
}
