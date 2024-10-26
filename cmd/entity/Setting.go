package entity

// swagger:model Setting
type Setting struct {
	// Setting ID
	// required: true
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// Description of the setting
	// required: true
	Description string `gorm:"type:varchar(255)"`
}
