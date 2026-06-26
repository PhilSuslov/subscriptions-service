package subscription

import (
	"context"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
)

func (uc *UseCase) TotalCost(ctx context.Context, f CostFilter) (int, error) {
	if f.PeriodTo.Before(f.PeriodFrom) {
		return 0, domain.ErrInvalidPeriod
	}
	
	return uc.repo.TotalCost(ctx, f)
}
