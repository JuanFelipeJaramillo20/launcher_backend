package entity

import "time"

// swagger:model Log
type Log struct {
	// Log ID
	// required: true
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// ID of the user associated with this log entry
	// required: true
	UserID uint64 `gorm:"index"`

	// Action performed
	// required: true
	Action string `gorm:"type:varchar(255)"`

	// Timestamp of the log entry
	// required: true
	Timestamp time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
