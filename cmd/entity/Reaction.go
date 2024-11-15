package entity

import "time"

// swagger:model Reaction
type Reaction struct {
	// Reaction ID
	// required: true
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// ID of the user who reacted to the news
	// required: true
	UserID uint64 `gorm:"index"`

	// ID of the news that was reacted to
	// required: true
	NewsID uint64 `gorm:"index"`

	// Type of reaction (e.g., "like" or "dislike")
	// required: true
	Type string `gorm:"type:varchar(10)"`

	// Timestamp of when the reaction occurred
	// required: true
	Timestamp time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
