package subscription

import (
	"fmt"
	"time"
)

const MonthLayout = "01-2006"

type Month struct {
	time.Time
}

func NewMonth(year int, month time.Month) (Month, error) {
	if year < 1 || month < time.January || month > time.December {
		return Month{}, fmt.Errorf("invalid month: year=%d month=%d", year, month)
	}
	return Month{Time: time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)}, nil
}

func ParseMonth(value string) (Month, error) {
	t, err := time.Parse(MonthLayout, value)
	if err != nil {
		return Month{}, fmt.Errorf("month must use MM-YYYY format: %w", err)
	}
	return Month{Time: time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)}, nil
}

func (m Month) String() string {
	return m.Time.Format(MonthLayout)
}

func (m Month) Before(other Month) bool { return m.Time.Before(other.Time) }
func (m Month) After(other Month) bool  { return m.Time.After(other.Time) }
