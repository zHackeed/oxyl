package models

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/oklog/ulid/v2"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`

	Name    string `json:"name"`
	Surname string `json:"surname"`
	Enabled bool   `json:"enabled"`

	MemberOf map[string]*CompanyMember `json:"member_of,omitempty"`

	LastLogin   time.Time `json:"last_login"`
	LastUpdated time.Time `json:"last_updated"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewUser(name, surname, email, hashedPassword string) (*User, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if surname == "" {
		return nil, fmt.Errorf("surname is required")
	}

	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	_, err := mail.ParseAddress(email) //https://stackoverflow.com/questions/66624011/how-to-validate-an-email-address-in-golang
	if err != nil {
		return nil, fmt.Errorf("invalid email address: %w", err)
	}

	return &User{
		ID:        ulid.Make().String(),
		Email:     email,
		Password:  hashedPassword,
		Name:      name,
		Surname:   surname,
		Enabled:   true,
		CreatedAt: time.Now(),
	}, nil
}

func (u *User) UpdateName(name, surname string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}

	if surname == "" {
		return fmt.Errorf("surname is required")
	}

	u.Name = name
	u.Surname = surname
	u.LastUpdated = time.Now()

	return nil
}

func (u *User) UpdatePassword(hashedPassword string) {
	u.Password = hashedPassword
	u.LastUpdated = time.Now()
}

func (u *User) UpdateEnabled(enabled bool) {
	u.Enabled = enabled
	u.LastUpdated = time.Now()
}

func (u *User) UpdateLastLogin() {
	u.LastLogin = time.Now()
}
