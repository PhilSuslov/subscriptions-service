package subscription

import (
	"context"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	"github.com/google/uuid"
)

type ListFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	Limit       int
	Offset      int
}

type CostFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	PeriodFrom  domain.Month
	PeriodTo    domain.Month
}

type Repository interface {
	Create(ctx context.Context, s *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, s *domain.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, f ListFilter) ([]domain.Subscription, error)
	TotalCost(ctx context.Context, f CostFilter) (int, error)
}
