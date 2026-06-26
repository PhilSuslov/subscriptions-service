package subscription

import (
	"context"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
)

func (uc *UseCase) List(ctx context.Context, f ListFilter) ([]domain.Subscription, error) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 50
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	return uc.repo.List(ctx, f)
}
