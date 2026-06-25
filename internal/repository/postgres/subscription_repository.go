package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	service "github.com/example/subscriptions-service/internal/service/subscription"
	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{pool: pool}
}

func (r *SubscriptionRepository) Create(ctx context.Context, s *domain.Subscription) error {
	query := `INSERT INTO subscriptions (id, service_name, price, user_id, start_month, end_month) VALUES ($1,$2,$3,$4,$5,$6) RETURNING created_at, updated_at`
	var end any
	if s.EndMonth != nil {
		end = s.EndMonth.Time
	}
	return r.pool.QueryRow(ctx, query, s.ID, s.ServiceName, s.Price, s.UserID, s.StartMonth.Time, end).Scan(&s.CreatedAt, &s.UpdatedAt)
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_month, end_month, created_at, updated_at FROM subscriptions WHERE id=$1`
	s, err := scanSubscription(r.pool.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return s, err
}

func (r *SubscriptionRepository) Update(ctx context.Context, s *domain.Subscription) error {
	query := `UPDATE subscriptions SET service_name=$2, price=$3, user_id=$4, start_month=$5, end_month=$6, updated_at=now() WHERE id=$1 RETURNING created_at, updated_at`
	var end any
	if s.EndMonth != nil {
		end = s.EndMonth.Time
	}
	err := r.pool.QueryRow(ctx, query, s.ID, s.ServiceName, s.Price, s.UserID, s.StartMonth.Time, end).Scan(&s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ct, err := r.pool.Exec(ctx, `DELETE FROM subscriptions WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *SubscriptionRepository) List(ctx context.Context, f service.ListFilter) ([]domain.Subscription, error) {
	where, args := buildCommonFilters(f.UserID, f.ServiceName)
	args = append(args, f.Limit, f.Offset)
	query := fmt.Sprintf(`SELECT id, service_name, price, user_id, start_month, end_month, created_at, updated_at FROM subscriptions %s ORDER BY created_at DESC, id DESC LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args))
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.Subscription, 0)
	for rows.Next() {
		s, err := scanSubscription(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *s)
	}
	return items, rows.Err()
}

func (r *SubscriptionRepository) TotalCost(ctx context.Context, f service.CostFilter) (int, error) {
	where, args := buildCommonFilters(f.UserID, f.ServiceName)
	periodArgsStart := len(args) + 1
	args = append(args, f.PeriodFrom.Time, f.PeriodTo.Time)
	if where == "" {
		where = "WHERE "
	} else {
		where += " AND "
	}
	where += fmt.Sprintf(`start_month <= $%d AND (end_month IS NULL OR end_month >= $%d)`, periodArgsStart+1, periodArgsStart)
	query := fmt.Sprintf(`
		SELECT COALESCE(SUM(price * (
			(EXTRACT(YEAR FROM age(LEAST(COALESCE(end_month, $%[2]d), $%[2]d), GREATEST(start_month, $%[1]d)))::int * 12) +
			EXTRACT(MONTH FROM age(LEAST(COALESCE(end_month, $%[2]d), $%[2]d), GREATEST(start_month, $%[1]d)))::int + 1
		)), 0)::int
		FROM subscriptions %[3]s`, periodArgsStart, periodArgsStart+1, where)
	var total int
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func buildCommonFilters(userID *uuid.UUID, serviceName *string) (string, []any) {
	clauses := make([]string, 0, 2)
	args := make([]any, 0, 2)
	if userID != nil {
		args = append(args, *userID)
		clauses = append(clauses, fmt.Sprintf("user_id=$%d", len(args)))
	}
	if serviceName != nil && strings.TrimSpace(*serviceName) != "" {
		args = append(args, strings.TrimSpace(*serviceName))
		clauses = append(clauses, fmt.Sprintf("service_name ILIKE $%d", len(args)))
	}
	if len(clauses) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(clauses, " AND "), args
}

func BuildCommonFiltersForTest(userID *uuid.UUID, serviceName *string) (string, []any) {
	return buildCommonFilters(userID, serviceName)
}

type scanner interface{ Scan(dest ...any) error }

func scanSubscription(row scanner) (*domain.Subscription, error) {
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
