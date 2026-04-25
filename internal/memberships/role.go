package memberships

const (
	RoleOwner = "owner"
	RoleAdmin = "admin"
	RoleStaff = "staff"
)

func CanAccessOperational(role string) bool { return true }
func CanAccessBilling(role string) bool { return role == RoleOwner || role == RoleAdmin }
func CanAccessUsers(role string) bool { return role == RoleOwner || role == RoleAdmin }
