package subscription

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartMonth  Month
	EndMonth    *Month
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func New(serviceName string, price int, userID uuid.UUID, start Month, end *Month) (*Subscription, error) {
	s := &Subscription{
		ID:          uuid.New(),
		ServiceName: strings.TrimSpace(serviceName),
		Price:       price,
		UserID:      userID,
		StartMonth:  start,
		EndMonth:    end,
	}
	if err := s.Validate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Subscription) Validate() error {
	if strings.TrimSpace(s.ServiceName) == "" {
		return ErrInvalidServiceName
	}
	if s.Price <= 0 {
		return ErrInvalidPrice
	}
	if s.UserID == uuid.Nil {
		return ErrInvalidUserID
	}
	if s.EndMonth != nil && s.EndMonth.Before(s.StartMonth) {
		return ErrInvalidPeriod
	}
	return nil
}
