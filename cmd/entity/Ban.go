package entity

import "time"

// swagger:model Ban
type Ban struct {
	// ID of the ban
	// required: true
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// Player ID associated with the ban
	// required: true
	PlayerID uint64 `gorm:"index"`

	// Reason for the ban
	// required: true
	Reason string `gorm:"type:varchar(255)"`

	// ID of the user who issued the ban
	// required: true
	BannedBy uint64 `gorm:"index"`

	// Date the ban was issued
	// required: true
	BanDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	// Duration of the ban
	Duration time.Duration
}
