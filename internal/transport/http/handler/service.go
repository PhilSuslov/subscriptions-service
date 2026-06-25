package handler

import (
	"context"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	app "github.com/example/subscriptions-service/internal/service/subscription"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	Create(ctx context.Context, cmd app.CreateCommand) (*domain.Subscription, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, cmd app.UpdateCommand) (*domain.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, f app.ListFilter) ([]domain.Subscription, error)
	TotalCost(ctx context.Context, f app.CostFilter) (int, error)
}
