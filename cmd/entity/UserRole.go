package entity

type UserRole struct {
	ID     uint64 `gorm:"primaryKey;autoIncrement"`
	UserID uint64 `gorm:"index"`
	RoleID uint64 `gorm:"index"`
}
