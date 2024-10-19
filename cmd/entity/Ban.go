package entity

import "time"

type Ban struct {
	ID       uint64    `gorm:"primaryKey;autoIncrement"`
	PlayerID uint64    `gorm:"index"`
	Reason   string    `gorm:"type:varchar(255)"`
	BannedBy uint64    `gorm:"index"`
	BanDate  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Duration time.Duration
}
