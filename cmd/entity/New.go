package entity

import "time"

type News struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Title     string    `gorm:"type:varchar(255)"`
	Content   string    `gorm:"type:text"`
	CreatedBy uint64    `gorm:"index"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
