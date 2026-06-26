package subscription

import domain "github.com/example/subscriptions-service/internal/domain/subscription"

type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
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
