package myapp

import "time"

// define user struct
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// define user struct to update
type UpdatedUser struct {
	ID               int    `json:"id"`
	UpdatedFirstName bool   `json:"updated_first_name"`
	FirstName        string `json:"first_name"`

	UpdatedLastName bool   `json:"updated_last_name"`
	LastName        string `json:"last_name"`

	UpdatedEmail bool   `json:"updated_email"`
	Email        string `json:"email"`

	UpdatedCreateAt bool      `json:"updated_created_at"`
	CreatedAt       time.Time `json:"created_at"`
}
