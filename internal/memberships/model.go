package memberships

import "time"

type Membership struct {
	TenantID string
	TenantSlug string
	TenantName string
	UserID string
	Role string
	CreatedAt time.Time
	UpdatedAt time.Time
}
