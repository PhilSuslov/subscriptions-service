package subscription

import (
	"context"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
)

func (uc *UseCase) Update(ctx context.Context, cmd UpdateCommand) (*domain.Subscription, error) {
	start, end, err := parsePeriod(cmd.StartDate, cmd.EndDate)
	if err != nil {
		return nil, err
	}

	s := &domain.Subscription{ID: cmd.ID, ServiceName: cmd.ServiceName, Price: cmd.Price,
		UserID: cmd.UserID, StartMonth: start, EndMonth: end}
	if err := s.Validate(); err != nil {
		return nil, err
	}

	return s, uc.repo.Update(ctx, s)
}
