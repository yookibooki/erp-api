package tenants

import "time"

type Tenant struct {
	ID string
	Slug string
	Name string
	Country string
	BaseCurrency string
	CreatedAt time.Time
	UpdatedAt time.Time
}
