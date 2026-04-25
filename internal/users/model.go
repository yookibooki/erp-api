package users

import "time"

type AppUser struct {
	ID string
	Email string
	DisplayName string
	CreatedAt time.Time
	UpdatedAt time.Time
}
