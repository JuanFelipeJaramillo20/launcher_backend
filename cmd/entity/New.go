package entity

import "time"

// swagger:model News
type News struct {
	// News ID
	// required: true
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// Title of the news article
	// required: true
	Title string `gorm:"type:varchar(255)"`

	// Content of the news article
	// required: true
	Content string `gorm:"type:text"`

	// ID of the user who created this news article
	// required: true
	CreatedBy uint64 `gorm:"index"`

	// Creation timestamp
	// required: true
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
