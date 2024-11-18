package utils

import (
	"fmt"

	"github.com/gorcon/rcon"
)

func NewRCONClient(address string, password string) (*rcon.Conn, error) {
	client, err := rcon.Dial(address, password)
	if err != nil {
		return nil, fmt.Errorf("error al conectar al servidor RCON: %w", err)
	}

	return client, nil
}
