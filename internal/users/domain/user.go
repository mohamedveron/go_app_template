package domain

import (
	"errors"
	"strings"
	"time"
)

// User holds all data required to represent a user
type User struct {
	FirstName string     `json:"firstName,omitempty"`
	LastName  string     `json:"lastName,omitempty"`
	Mobile    string     `json:"mobile,omitempty"`
	Email     string     `json:"email,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

func (u *User) SetDefaults() {
	now := time.Now()
	if u.CreatedAt == nil {
		u.CreatedAt = &now
	}

	if u.UpdatedAt == nil {
		u.UpdatedAt = &now
	}
}

// Sanitize is used to sanitize/cleanup the fields of User
func (u *User) Sanitize() {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(u.Email)
	u.Mobile = strings.TrimSpace(u.Mobile)
}

// Validate is used to validate the fields of User
func (u *User) Validate() error {
	if u.Email == "" {
		return nil
	}

	err := u.ValidateEmail(u.Email)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) ValidateEmail(email string) error {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("invalid email address provided")
	}

	return nil
}
