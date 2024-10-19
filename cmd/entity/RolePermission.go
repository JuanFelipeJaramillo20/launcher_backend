package entity

type RolePermission struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	RoleID       uint64 `gorm:"index"`
	PermissionID uint64 `gorm:"index"`
}
