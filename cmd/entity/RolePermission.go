package entity

// swagger:model RolePermission
type RolePermission struct {
	// RolePermission ID
	// required: true
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// Role ID
	// required: true
	RoleID uint64 `gorm:"index"`

	// Permission ID
	// required: true
	PermissionID uint64 `gorm:"index"`
}
