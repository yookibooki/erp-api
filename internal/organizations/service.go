package organizations

import "context"

type Service struct{}

type CreateInput struct {
	OwnerUserID string
	CompanyName string
	Slug string
	Country string
	BaseCurrency string
}

func NewService() *Service { return &Service{} }
func (s *Service) Create(ctx context.Context, in CreateInput) (any, error) { return nil, nil }
