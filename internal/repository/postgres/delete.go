package postgres

import (
	"context"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	"github.com/google/uuid"
)

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
