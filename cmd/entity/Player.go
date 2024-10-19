package entity

import "time"

type Player struct {
	ID       uint64    `gorm:"primaryKey;autoIncrement"`
	UserID   uint64    `gorm:"index"`
	ServerID uint64    `gorm:"index"`
	JoinDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	LastSeen time.Time
}
