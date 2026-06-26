package subscription

import (
	"context"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
)

func (uc *UseCase) Create(ctx context.Context, cmd CreateCommand) (*domain.Subscription, error) {
	start, end, err := parsePeriod(cmd.StartDate, cmd.EndDate)
	if err != nil {
		return nil, err
	}

	s, err := domain.New(cmd.ServiceName, cmd.Price, cmd.UserID, start, end)
	if err != nil {
		return nil, err
	}
	
	return s, uc.repo.Create(ctx, s)
}
