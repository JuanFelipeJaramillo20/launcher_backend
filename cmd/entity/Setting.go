package entity

type Setting struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	Description string `gorm:"type:varchar(255)"`
}
