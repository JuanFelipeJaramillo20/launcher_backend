package utils

import (
	"errors"
	"fmt"
	"regexp"
)

const MinPasswordLength = 8

func IsValidEmail(email string) bool {
	const emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func IsValidNickname(nickname string) bool {
	const nicknameRegex = `^[a-zA-Z0-9_]{3,30}$`
	re := regexp.MustCompile(nicknameRegex)
	return re.MatchString(nickname)
}

func IsValidPassword(password string) error {
	if len(password) < MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", MinPasswordLength)
	}

	const (
		upperCase = `[A-Z]`
		digit     = `[0-9]`
		special   = `[\W_]`
	)

	if match, _ := regexp.MatchString(upperCase, password); !match {
		return errors.New("password must contain at least one uppercase letter")
	}
	if match, _ := regexp.MatchString(digit, password); !match {
		return errors.New("password must contain at least one digit")
	}
	if match, _ := regexp.MatchString(special, password); !match {
		return errors.New("password must contain at least one special character")
	}
	return nil
}
