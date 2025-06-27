package models

import "time"

type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Password   string    `json:"-"` // Omit from JSON responses
	AuthMethod string    `json:"auth_method"`
	CreatedAt  time.Time `json:"created_at"`
}
