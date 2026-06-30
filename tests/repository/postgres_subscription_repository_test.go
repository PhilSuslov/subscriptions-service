package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	postgres "github.com/example/subscriptions-service/internal/repository/postgres"
	service "github.com/example/subscriptions-service/internal/service/subscription"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

type row struct {
	scan func(dest ...interface{}) error
}

func (r row) Scan(dest ...interface{}) error { return r.scan(dest...) }

type rows struct {
	steps int
	err   error
}

func (r *rows) Close()                                       {}
func (r *rows) Err() error                                   { return r.err }
func (r *rows) CommandTag() pgconn.CommandTag                { return pgconn.NewCommandTag("SELECT 0") }
func (r *rows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (r *rows) Next() bool                                   { r.steps++; return r.steps == 1 && r.err == nil }
func (r *rows) Scan(dest ...interface{}) error               { return nil }
func (r *rows) Values() ([]interface{}, error)               { return nil, nil }
func (r *rows) RawValues() [][]byte                          { return nil }
func (r *rows) Conn() *pgx.Conn                              { return nil }

type queryer struct {
	rowFn func(sql string, args ...interface{}) pgx.Row
	qFn   func(sql string, args ...interface{}) (pgx.Rows, error)
	eFn   func(sql string, args ...interface{}) (pgconn.CommandTag, error)
}

func (q queryer) QueryRow(_ context.Context, sql string, args ...interface{}) pgx.Row {
	return q.rowFn(sql, args...)
}
func (q queryer) Query(_ context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return q.qFn(sql, args...)
}
func (q queryer) Exec(_ context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return q.eFn(sql, args...)
}

func TestBuildCommonFilters(t *testing.T) {
	userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
	serviceName := " Yandex Plus "
	where, args := postgres.BuildCommonFiltersForTest(&userID, &serviceName)
	require.Equal(t, "WHERE user_id=$1 AND service_name ILIKE $2", where)
	require.Len(t, args, 2)
}

func TestRepositoryHappyPath(t *testing.T) {
	id := uuid.New()
	start := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	now := time.Date(2025, 6, 25, 12, 0, 0, 0, time.UTC)
	q := queryer{
		rowFn: func(sql string, args ...interface{}) pgx.Row {
			switch {
			case sql[:6] == "INSERT":
				return row{scan: func(dest ...interface{}) error {
					*(dest[0].(*time.Time)) = now
					*(dest[1].(*time.Time)) = now
					return nil
				}}
			case sql[:6] == "UPDATE":
				return row{scan: func(dest ...interface{}) error {
					*(dest[0].(*time.Time)) = now
					*(dest[1].(*time.Time)) = now
					return nil
				}}
			case sql[:6] == "SELECT":
				return row{scan: func(dest ...interface{}) error {
					*(dest[0].(*uuid.UUID)) = id
					*(dest[1].(*string)) = "Yandex Plus"
					*(dest[2].(*int)) = 400
					*(dest[3].(*uuid.UUID)) = uuid.New()
					*(dest[4].(*time.Time)) = start
					*(dest[5].(**time.Time)) = &end
					*(dest[6].(*time.Time)) = now
					*(dest[7].(*time.Time)) = now
					return nil
				}}
			}
			return row{scan: func(dest ...interface{}) error { return nil }}
		},
		qFn: func(string, ...interface{}) (pgx.Rows, error) { return &rows{}, nil },
		eFn: func(string, ...interface{}) (pgconn.CommandTag, error) { return pgconn.NewCommandTag("DELETE 1"), nil },
	}
	repo := postgres.NewSubscriptionRepository(q)
	s := &domain.Subscription{ID: id, ServiceName: "Yandex Plus", Price: 400, UserID: uuid.New(), StartMonth: domain.Month{Time: start}, EndMonth: &domain.Month{Time: end}}
	require.NoError(t, repo.Create(context.Background(), s))
	_, _ = repo.GetByID(context.Background(), id)
	require.NoError(t, repo.Update(context.Background(), s))
	require.NoError(t, repo.Delete(context.Background(), id))
	_, _ = repo.List(context.Background(), service.ListFilter{Limit: 10, Offset: 0})
	_, _ = repo.TotalCost(context.Background(), service.CostFilter{PeriodFrom: domain.Month{Time: start}, PeriodTo: domain.Month{Time: start}})
}

func TestRepositoryErrors(t *testing.T) {
	repo := postgres.NewSubscriptionRepository(queryer{
		rowFn: func(string, ...interface{}) pgx.Row {
			return row{scan: func(dest ...interface{}) error { return pgx.ErrNoRows }}
		},
		qFn: func(string, ...interface{}) (pgx.Rows, error) { return nil, errors.New("query") },
		eFn: func(string, ...interface{}) (pgconn.CommandTag, error) { return pgconn.NewCommandTag("DELETE 0"), nil },
	})
	_, err := repo.GetByID(context.Background(), uuid.New())
	require.ErrorIs(t, err, domain.ErrNotFound)
	err = repo.Delete(context.Background(), uuid.New())
	require.ErrorIs(t, err, domain.ErrNotFound)
	_, err = repo.List(context.Background(), service.ListFilter{})
	require.Error(t, err)
}

func TestRepositoryNilEndMonth(t *testing.T) {
	start := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	now := time.Date(2025, 6, 25, 12, 0, 0, 0, time.UTC)
	called := false
	repo := postgres.NewSubscriptionRepository(queryer{
		rowFn: func(sql string, args ...interface{}) pgx.Row {
			called = true
			return row{scan: func(dest ...interface{}) error {
				*(dest[0].(*time.Time)) = now
				*(dest[1].(*time.Time)) = now
				return nil
			}}
		},
	})
	s := &domain.Subscription{ID: uuid.New(), ServiceName: "A", Price: 1, UserID: uuid.New(), StartMonth: domain.Month{Time: start}}
	require.NoError(t, repo.Create(context.Background(), s))
	require.True(t, called)
}

func TestBuildCommonFiltersPartial(t *testing.T) {
	serviceName := "Yandex"
	where, args := postgres.BuildCommonFiltersForTest(nil, &serviceName)
	require.Equal(t, "WHERE service_name ILIKE $1", where)
	require.Len(t, args, 1)
}
