package entity

import "time"

type Log struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	UserID    uint64    `gorm:"index"`
	Action    string    `gorm:"type:varchar(255)"`
	Timestamp time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
