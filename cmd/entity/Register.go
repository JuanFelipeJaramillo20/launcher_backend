package entity

type Register struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement"`
	FullName string `gorm:"type:varchar(255)"`
	Email    string `gorm:"type:varchar(255);unique"`
	Nickname string `gorm:"type:varchar(100);unique"`
	Password string `gorm:"type:varchar(255)"`
	Accepted bool   `gorm:"default:false"`
}
