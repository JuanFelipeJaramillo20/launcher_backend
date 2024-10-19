package entity

type Role struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(100);unique"`
}
