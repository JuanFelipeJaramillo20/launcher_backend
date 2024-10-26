package entity

// swagger:model Permission
type Permission struct {
	// Permission ID
	// required: true
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// Name of the permission
	// required: true
	Name string `gorm:"type:varchar(100);unique"`
}
