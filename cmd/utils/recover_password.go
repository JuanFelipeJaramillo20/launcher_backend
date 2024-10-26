package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func GenerateResetToken() (string, time.Time, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", time.Time{}, err
	}
	token := hex.EncodeToString(bytes)
	expiration := time.Now().Add(15 * time.Minute)
	return token, expiration, nil
}
