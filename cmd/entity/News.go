package entity

import "time"

// swagger:model News
type News struct {
	// News ID
	// required: true
	ID uint64 `json:"id" gorm:"primaryKey;autoIncrement"`

	// Title of the news article
	// required: true
	Title string `json:"title" gorm:"type:varchar(255)"`

	// Content of the news article
	// required: true
	Content string `json:"content" gorm:"type:text"`

	// ID of the user who created this news article
	// required: true
	CreatedBy uint64 `json:"created_by" gorm:"index"`

	// Creation timestamp
	// required: true
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}
