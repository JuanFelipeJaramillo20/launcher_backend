package entity

import "time"

// swagger:model Like
type Like struct {
	// Like ID
	// required: true
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// ID of the user who liked the news
	// required: true
	UserID uint64 `gorm:"index"`

	// ID of the news that was liked
	// required: true
	NewsID uint64 `gorm:"index"`

	// Timestamp of when the like occurred
	// required: true
	Timestamp time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
