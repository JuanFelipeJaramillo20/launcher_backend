package entity

// swagger:model Role
type Role struct {
	// Role ID
	// required: true
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// Name of the role
	// required: true
	Name string `gorm:"type:varchar(100);unique"`
}
