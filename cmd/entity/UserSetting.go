package entity

type UserSetting struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	UserID    uint64 `gorm:"index"`
	SettingID uint64 `gorm:"index"`
	IsActive  bool   `gorm:"default:true"`
}
