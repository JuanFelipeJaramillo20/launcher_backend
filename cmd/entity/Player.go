package entity

import "time"

// swagger:model Player
type Player struct {
	// Player ID
	// required: true
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// Associated user ID
	// required: true
	UserID uint64 `gorm:"index"`

	// Server ID the player is associated with
	// required: true
	ServerID uint64 `gorm:"index"`

	// Date the player joined
	// required: true
	JoinDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	// Last time the player was seen
	LastSeen time.Time
}
