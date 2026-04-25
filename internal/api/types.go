package api

import "time"

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type ListMeta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type UserDTO struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}

type OrganizationDTO struct {
	Slug         string    `json:"slug"`
	Name         string    `json:"name"`
	Country      string    `json:"country"`
	BaseCurrency string    `json:"baseCurrency"`
	Role         string    `json:"role"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type PermissionsDTO struct {
	CanAccessOperational bool `json:"canAccessOperational"`
	CanManageBilling     bool `json:"canManageBilling"`
	CanManageMembers     bool `json:"canManageMembers"`
	CanManageSettings    bool `json:"canManageSettings"`
}

// MeResponse is pure API data - NO CSRF token, NO cookies
type MeResponse struct {
	User               UserDTO           `json:"user"`
	Organizations      []OrganizationDTO `json:"organizations"`
	ActiveOrganization *OrganizationDTO  `json:"activeOrganization,omitempty"`
	Permissions        PermissionsDTO    `json:"permissions"`
}

type CurrentOrganizationResponse struct {
	Organization *OrganizationDTO `json:"organization,omitempty"`
}

type CreateOrganizationRequest struct {
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Country      string `json:"country"`
	BaseCurrency string `json:"baseCurrency"`
}

type OrganizationListResponse struct {
	Items []OrganizationDTO `json:"items"`
	Meta  ListMeta          `json:"meta"`
}

type CreateOrganizationResponse struct {
	Organization OrganizationDTO `json:"organization"`
}

// ... rest of DTOs unchanged (CustomerDTO, ItemDTO, InvoiceDTO, etc.)
type CustomerDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     *string   `json:"email,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CustomerListResponse struct {
	Items []CustomerDTO `json:"items"`
	Meta  ListMeta      `json:"meta"`
}

type CreateCustomerRequest struct {
	Name  string  `json:"name"`
	Email *string `json:"email,omitempty"`
}

type UpdateCustomerRequest struct {
	Name  string  `json:"name"`
	Email *string `json:"email,omitempty"`
}

type ItemDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	UnitPrice string    `json:"unitPrice"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ItemListResponse struct {
	Items []ItemDTO `json:"items"`
	Meta  ListMeta  `json:"meta"`
}

type CreateItemRequest struct {
	Name      string `json:"name"`
	UnitPrice string `json:"unitPrice"`
}

type UpdateItemRequest struct {
	Name      string `json:"name"`
	UnitPrice string `json:"unitPrice"`
}

type InvoiceLineInput struct {
	ItemID string `json:"itemId"`
	Qty    int    `json:"qty"`
}

type InvoiceLineDTO struct {
	ID        string  `json:"id"`
	ItemID    *string `json:"itemId,omitempty"`
	ItemName  string  `json:"itemName"`
	Qty       int     `json:"qty"`
	UnitPrice string  `json:"unitPrice"`
	LineTotal string  `json:"lineTotal"`
}

type InvoiceSummaryDTO struct {
	ID                string     `json:"id"`
	CustomerID        string     `json:"customerId"`
	InvoiceNumber     string     `json:"invoiceNumber"`
	CustomerName      string     `json:"customerName"`
	TotalAmount       string     `json:"totalAmount"`
	PaidAmount        string     `json:"paidAmount"`
	OutstandingAmount string     `json:"outstandingAmount"`
	IssuedAt          time.Time  `json:"issuedAt"`
	SentAt            *time.Time `json:"sentAt,omitempty"`
	VoidedAt          *time.Time `json:"voidedAt,omitempty"`
	Status            string     `json:"status"`
}

type InvoiceListResponse struct {
	Items []InvoiceSummaryDTO `json:"items"`
	Meta  ListMeta            `json:"meta"`
}

type CreateInvoiceRequest struct {
	CustomerID string             `json:"customerId"`
	Lines      []InvoiceLineInput `json:"lines"`
}

type UpdateInvoiceRequest struct {
	CustomerID string             `json:"customerId"`
	Lines      []InvoiceLineInput `json:"lines"`
}

type InvoiceDetailDTO struct {
	Invoice  InvoiceSummaryDTO `json:"invoice"`
	Customer CustomerDTO       `json:"customer"`
	Lines    []InvoiceLineDTO  `json:"lines"`
	Payments []PaymentDTO      `json:"payments"`
}

type PaymentDTO struct {
	ID            string    `json:"id"`
	InvoiceID     string    `json:"invoiceId"`
	InvoiceNumber string    `json:"invoiceNumber"`
	Amount        string    `json:"amount"`
	PaidAt        time.Time `json:"paidAt"`
}

type PaymentListResponse struct {
	Items []PaymentDTO `json:"items"`
	Meta  ListMeta     `json:"meta"`
}

type CreatePaymentRequest struct {
	InvoiceID string `json:"invoiceId"`
	Amount    string `json:"amount"`
}

type ReportTotalsDTO struct {
	InvoiceCount     int    `json:"invoiceCount"`
	PaymentCount     int    `json:"paymentCount"`
	TotalInvoiced    string `json:"totalInvoiced"`
	TotalPaid        string `json:"totalPaid"`
	TotalOutstanding string `json:"totalOutstanding"`
}

type ReportSummaryResponse struct {
	Summary        ReportTotalsDTO     `json:"summary"`
	RecentInvoices []InvoiceSummaryDTO `json:"recentInvoices"`
	RecentPayments []PaymentDTO        `json:"recentPayments"`
}

type MembershipDTO struct {
	OrganizationSlug         string    `json:"organizationSlug"`
	OrganizationName         string    `json:"organizationName"`
	OrganizationCountry      string    `json:"organizationCountry"`
	OrganizationBaseCurrency string    `json:"organizationBaseCurrency"`
	UserID                   string    `json:"userId"`
	UserEmail                string    `json:"userEmail"`
	UserDisplayName          string    `json:"userDisplayName"`
	Role                     string    `json:"role"`
	IsCurrent                bool      `json:"isCurrent"`
	CreatedAt                time.Time `json:"createdAt"`
	UpdatedAt                time.Time `json:"updatedAt"`
}

type MembershipListResponse struct {
	Items []MembershipDTO `json:"items"`
	Meta  ListMeta        `json:"meta"`
}

type UpdateMembershipRequest struct {
	Role string `json:"role"`
}

type InvitationDTO struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"createdAt"`
	ExpiresAt  time.Time  `json:"expiresAt"`
	AcceptedAt *time.Time `json:"acceptedAt,omitempty"`
	RevokedAt  *time.Time `json:"revokedAt,omitempty"`
}

type InvitationListResponse struct {
	Items []InvitationDTO `json:"items"`
	Meta  ListMeta        `json:"meta"`
}

type CreateInvitationRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type CreateInvitationResponse struct {
	Invitation InvitationDTO `json:"invitation"`
	Token      string        `json:"token"`
	InviteURL  string        `json:"inviteUrl"`
}

type AcceptInvitationRequest struct {
	Token string `json:"token"`
}

type NumberingConfigDTO struct {
	Prefix     string `json:"prefix"`
	NextNumber int64  `json:"nextNumber"`
	Padding    int    `json:"padding"`
}

type SettingsDTO struct {
	Organization OrganizationDTO    `json:"organization"`
	Numbering    NumberingConfigDTO `json:"numbering"`
}

type UpdateSettingsRequest struct {
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Country      string `json:"country"`
	BaseCurrency string `json:"baseCurrency"`
	Prefix       string `json:"prefix"`
	NextNumber   int64  `json:"nextNumber"`
	Padding      int    `json:"padding"`
}
