package subscription

import (
	"context"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	"github.com/google/uuid"
)

func (uc *UseCase) Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return uc.repo.GetByID(ctx, id)
}
