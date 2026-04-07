package entity

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Role     string    `json:"role"`
	Verified bool      `json:"verified"`
}
