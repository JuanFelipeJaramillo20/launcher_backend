package entity

type Server struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement"`
	Name       string `gorm:"type:varchar(100)"`
	IPAddress  string `gorm:"type:varchar(45)"`
	Port       int    `gorm:"type:int"`
	MaxPlayers int    `gorm:"type:int"`
	Status     string `gorm:"type:varchar(50)"`
}
