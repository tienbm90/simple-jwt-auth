package models

import (
	"fmt"
)

type User struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Todo struct {
	UserID string `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
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
