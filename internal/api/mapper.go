package api

import (
	"time"

	"github.com/yookibooki/erp-api/internal/app"
	"github.com/yookibooki/erp-api/internal/customers"
	"github.com/yookibooki/erp-api/internal/invitations"
	"github.com/yookibooki/erp-api/internal/invoices"
	"github.com/yookibooki/erp-api/internal/items"
	"github.com/yookibooki/erp-api/internal/memberships"
	"github.com/yookibooki/erp-api/internal/numbering"
	"github.com/yookibooki/erp-api/internal/payments"
	"github.com/yookibooki/erp-api/internal/reports"
	"github.com/yookibooki/erp-api/internal/tenants"
	"github.com/yookibooki/erp-api/internal/users"
)

func mapUser(u users.AppUser) UserDTO {
	return UserDTO{
		ID:          u.ID,
		Email:       u.Email,
		DisplayName: u.DisplayName,
	}
}

func mapOrganizationFromMembership(m memberships.Membership, active bool) OrganizationDTO {
	return OrganizationDTO{
		Slug:         m.TenantSlug,
		Name:         m.TenantName,
		Country:      m.TenantCountry,
		BaseCurrency: m.TenantBaseCurrency,
		Role:         m.Role,
		Active:       active,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func mapOrganization(t tenants.Tenant, role string, active bool) OrganizationDTO {
	return OrganizationDTO{
		Slug:         t.Slug,
		Name:         t.Name,
		Country:      t.Country,
		BaseCurrency: t.BaseCurrency,
		Role:         role,
		Active:       active,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}

func mapPermissions(role string) PermissionsDTO {
	return PermissionsDTO{
		CanAccessOperational: memberships.CanAccessOperational(role),
		CanManageBilling:     memberships.CanAccessBilling(role),
		CanManageMembers:     memberships.CanAccessUsers(role),
		CanManageSettings:    role == memberships.RoleOwner || role == memberships.RoleAdmin,
	}
}

// mapMe is now PURE - no CSRF, no session
func mapMe(identity app.RequestContext, membershipsList []memberships.Membership) MeResponse {
	user := UserDTO{
		ID:          identity.UserID,
		Email:       identity.Email,
		DisplayName: identity.DisplayName,
	}

	orgs := make([]OrganizationDTO, 0, len(membershipsList))
	var active *OrganizationDTO
	var perms PermissionsDTO

	for _, m := range membershipsList {
		isActive := m.TenantID == identity.OrganizationIDOrTenant()
		dto := mapOrganizationFromMembership(m, isActive)
		orgs = append(orgs, dto)
		if isActive {
			active = &dto
			perms = mapPermissions(m.Role)
		}
	}

	return MeResponse{
		User:               user,
		Organizations:      orgs,
		ActiveOrganization: active,
		Permissions:        perms,
	}
}

func mapCustomer(c customers.Customer) CustomerDTO {
	return CustomerDTO{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func mapItem(i items.Item) ItemDTO {
	return ItemDTO{
		ID:        i.ID,
		Name:      i.Name,
		UnitPrice: i.UnitPrice,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}

func mapInvoiceSummary(v invoices.InvoiceSummary) InvoiceSummaryDTO {
	return InvoiceSummaryDTO{
		ID:                v.ID,
		CustomerID:        v.CustomerID,
		InvoiceNumber:     v.InvoiceNumber,
		CustomerName:      v.CustomerName,
		TotalAmount:       v.TotalAmount,
		PaidAmount:        v.PaidAmount,
		OutstandingAmount: v.OutstandingAmount,
		IssuedAt:          v.IssuedAt,
		SentAt:            v.SentAt,
		VoidedAt:          v.VoidedAt,
		Status:            v.Status,
	}
}

func mapInvoiceLine(v invoices.InvoiceItem) InvoiceLineDTO {
	return InvoiceLineDTO{
		ID:        v.ID,
		ItemID:    v.ItemID,
		ItemName:  v.ItemName,
		Qty:       v.Qty,
		UnitPrice: v.UnitPrice,
		LineTotal: v.LineTotal,
	}
}

func mapPayment(v payments.PaymentSummary) PaymentDTO {
	return PaymentDTO{
		ID:            v.ID,
		InvoiceID:     v.InvoiceID,
		InvoiceNumber: v.InvoiceNumber,
		Amount:        v.Amount,
		PaidAt:        v.PaidAt,
	}
}

func mapReportTotals(v reports.Summary) ReportTotalsDTO {
	return ReportTotalsDTO{
		InvoiceCount:     v.InvoiceCount,
		PaymentCount:     v.PaymentCount,
		TotalInvoiced:    v.TotalInvoiced,
		TotalPaid:        v.TotalPaid,
		TotalOutstanding: v.TotalOutstanding,
	}
}

func mapMembership(m memberships.Membership, current bool) MembershipDTO {
	return MembershipDTO{
		OrganizationSlug:         m.TenantSlug,
		OrganizationName:         m.TenantName,
		OrganizationCountry:      m.TenantCountry,
		OrganizationBaseCurrency: m.TenantBaseCurrency,
		UserID:                   m.UserID,
		UserEmail:                m.UserEmail,
		UserDisplayName:          m.UserDisplayName,
		Role:                     m.Role,
		IsCurrent:                current,
		CreatedAt:                m.CreatedAt,
		UpdatedAt:                m.UpdatedAt,
	}
}

func mapInvite(inv invitations.Invite) InvitationDTO {
	return InvitationDTO{
		ID:         inv.ID,
		Email:      inv.Email,
		Role:       inv.Role,
		Status:     inviteStatus(inv),
		CreatedAt:  inv.CreatedAt,
		ExpiresAt:  inv.ExpiresAt,
		AcceptedAt: inv.AcceptedAt,
		RevokedAt:  inv.RevokedAt,
	}
}

func mapNumbering(cfg numbering.Config) NumberingConfigDTO {
	return NumberingConfigDTO{
		Prefix:     cfg.Prefix,
		NextNumber: cfg.NextNumber,
		Padding:    cfg.Padding,
	}
}

func inviteStatus(inv invitations.Invite) string {
	if inv.AcceptedAt != nil {
		return "accepted"
	}
	if inv.RevokedAt != nil {
		return "revoked"
	}
	if timeNow().After(inv.ExpiresAt) {
		return "expired"
	}
	return "pending"
}

var timeNow = func() time.Time { return time.Now().UTC() }
