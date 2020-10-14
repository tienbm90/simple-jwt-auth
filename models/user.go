package models

import (
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName      string `json:"username"`
	Password      string `json:"password"`
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email" gorm:"unique"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

// SetPassword sets a new password stored as hash.
func (m *User) SetPassword(password string) error {

	if len(password) < 6 {
		return fmt.Errorf("new password for %s must be at least 6 characters", m.UserName)
	}
	m.Password = password
	return nil
}

// InvalidPassword returns true if the given password does not match the hash.
func (m *User) InvalidPassword(password string) bool {

	if password == "" {
		return true
	}

	if m.Password != password {
		return true
	}

	return false
}
