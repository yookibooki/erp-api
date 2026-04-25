package app

type RequestContext struct {
	UserID           string
	Email            string
	DisplayName      string
	OrganizationID   string
	OrganizationSlug string
}

func (rc RequestContext) OrganizationIDOrTenant() string { return rc.OrganizationID }
func (rc RequestContext) OrganizationSlugOrTenant() string { return rc.OrganizationSlug }
