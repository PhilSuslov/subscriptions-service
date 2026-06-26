package postgres

import (
	"time"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
)

func scanSubscription(row interface{ Scan(dest ...any) error }) (*domain.Subscription, error) {
	var s domain.Subscription
	var start time.Time
	var end *time.Time
	if err := row.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &start, &end, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return nil, err
	}
	s.StartMonth = domain.Month{Time: start}
	if end != nil {
		m := domain.Month{Time: *end}
		s.EndMonth = &m
	}
	return &s, nil
}
