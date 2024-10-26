package entity

// swagger:model UserRole
type UserRole struct {
	// UserRole ID
	// required: true
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// User ID
	// required: true
	UserID uint64 `gorm:"index"`

	// Role ID
	// required: true
	RoleID uint64 `gorm:"index"`
}
