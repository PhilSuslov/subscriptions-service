package subscription

import (
	"context"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	"github.com/google/uuid"
)

type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

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

func (uc *UseCase) Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) Update(ctx context.Context, cmd UpdateCommand) (*domain.Subscription, error) {
	start, end, err := parsePeriod(cmd.StartDate, cmd.EndDate)
	if err != nil {
		return nil, err
	}
	s := &domain.Subscription{ID: cmd.ID, ServiceName: cmd.ServiceName, Price: cmd.Price, UserID: cmd.UserID, StartMonth: start, EndMonth: end}
	if err := s.Validate(); err != nil {
		return nil, err
	}
	return s, uc.repo.Update(ctx, s)
}

func (uc *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *UseCase) List(ctx context.Context, f ListFilter) ([]domain.Subscription, error) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 50
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	return uc.repo.List(ctx, f)
}

func (uc *UseCase) TotalCost(ctx context.Context, f CostFilter) (int, error) {
	if f.PeriodTo.Before(f.PeriodFrom) {
		return 0, domain.ErrInvalidPeriod
	}
	return uc.repo.TotalCost(ctx, f)
}

func parsePeriod(startDate string, endDate *string) (domain.Month, *domain.Month, error) {
	start, err := domain.ParseMonth(startDate)
	if err != nil {
		return domain.Month{}, nil, err
	}
	if endDate == nil || *endDate == "" {
		return start, nil, nil
	}
	end, err := domain.ParseMonth(*endDate)
	if err != nil {
		return domain.Month{}, nil, err
	}
	return start, &end, nil
}
