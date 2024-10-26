package entity

// swagger:model UserSetting
type UserSetting struct {
	// UserSetting ID
	// required: true
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// User ID
	// required: true
	UserID uint64 `gorm:"index"`

	// Setting ID
	// required: true
	SettingID uint64 `gorm:"index"`

	// Active status of the setting
	// required: true
	IsActive bool `gorm:"default:true"`
}
