package entity

// swagger:model Server
type Server struct {
	// Server ID
	// required: true
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	// Name of the server
	// required: true
	Name string `gorm:"type:varchar(100)"`

	// IP address of the server
	// required: true
	IPAddress string `gorm:"type:varchar(45)"`

	// Port number of the server
	// required: true
	Port int `gorm:"type:int"`

	// Maximum number of players on the server
	// required: true
	MaxPlayers int `gorm:"type:int"`

	// Current status of the server
	// required: true
	Status string `gorm:"type:varchar(50)"`
}
