package service_test

import (
	"context"
	"sync"
	"testing"

	domain "github.com/example/subscriptions-service/internal/domain/subscription"
	service "github.com/example/subscriptions-service/internal/service/subscription"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type fakeRepo struct {
	mu    sync.Mutex
	items map[uuid.UUID]*domain.Subscription
}

func newFakeRepo() *fakeRepo { return &fakeRepo{items: make(map[uuid.UUID]*domain.Subscription)} }

func (r *fakeRepo) Create(_ context.Context, s *domain.Subscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[s.ID] = s
	return nil
}
func (r *fakeRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Subscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return s, nil
}
func (r *fakeRepo) Update(_ context.Context, s *domain.Subscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[s.ID]; !ok {
		return nil
	}
	r.items[s.ID] = s
	return nil
}
func (r *fakeRepo) Delete(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.items, id)
	return nil
}
func (r *fakeRepo) List(_ context.Context, _ service.ListFilter) ([]domain.Subscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]domain.Subscription, 0, len(r.items))
	for _, s := range r.items {
		out = append(out, *s)
	}
	return out, nil
}
func (r *fakeRepo) TotalCost(_ context.Context, f service.CostFilter) (int, error) {
	if f.PeriodTo.Before(f.PeriodFrom) {
		return 0, domain.ErrInvalidPeriod
	}
	return 400, nil
}

func TestUseCaseCreate(t *testing.T) {
	uc := service.NewUseCase(newFakeRepo())
	s, err := uc.Create(context.Background(), service.CreateCommand{ServiceName: "Yandex Plus", Price: 400, UserID: uuid.New(), StartDate: "07-2025"})
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, s.ID)
}

func TestUseCaseCreateInvalidDate(t *testing.T) {
	uc := service.NewUseCase(newFakeRepo())
	_, err := uc.Create(context.Background(), service.CreateCommand{ServiceName: "Yandex Plus", Price: 400, UserID: uuid.New(), StartDate: "2025-07"})
	require.Error(t, err)
}

func TestUseCaseCreateRejectsInvalidFields(t *testing.T) {
	uc := service.NewUseCase(newFakeRepo())

	_, err := uc.Create(context.Background(), service.CreateCommand{ServiceName: "   ", Price: 400, UserID: uuid.New(), StartDate: "07-2025"})
	require.ErrorIs(t, err, domain.ErrInvalidServiceName)

	_, err = uc.Create(context.Background(), service.CreateCommand{ServiceName: "Yandex Plus", Price: 0, UserID: uuid.New(), StartDate: "07-2025"})
	require.ErrorIs(t, err, domain.ErrInvalidPrice)

	_, err = uc.Create(context.Background(), service.CreateCommand{ServiceName: "Yandex Plus", Price: 400, UserID: uuid.Nil, StartDate: "07-2025"})
	require.ErrorIs(t, err, domain.ErrInvalidUserID)
}

func TestUseCaseListNormalizesPagination(t *testing.T) {
	repo := newFakeRepo()
	uc := service.NewUseCase(repo)

	items, err := uc.List(context.Background(), service.ListFilter{Limit: 0, Offset: -10})
	require.NoError(t, err)
	require.Empty(t, items)
}

func TestUseCaseListAcceptsUpperBoundLimit(t *testing.T) {
	repo := newFakeRepo()
	uc := service.NewUseCase(repo)

	_, err := uc.List(context.Background(), service.ListFilter{Limit: 100, Offset: 0})
	require.NoError(t, err)
}

func TestUseCaseDeleteMissingReturnsNotFound(t *testing.T) {
	uc := service.NewUseCase(newFakeRepo())
	err := uc.Delete(context.Background(), uuid.New())
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestUseCaseGetMissingReturnsNotFound(t *testing.T) {
	uc := service.NewUseCase(newFakeRepo())
	_, err := uc.Get(context.Background(), uuid.New())
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestUseCaseTotalCostInvalidPeriod(t *testing.T) {
	uc := service.NewUseCase(newFakeRepo())
	from, _ := domain.ParseMonth("12-2025")
	to, _ := domain.ParseMonth("07-2025")
	_, err := uc.TotalCost(context.Background(), service.CostFilter{PeriodFrom: from, PeriodTo: to})
	require.ErrorIs(t, err, domain.ErrInvalidPeriod)
}

func TestUseCaseUpdateInvalidPeriod(t *testing.T) {
	uc := service.NewUseCase(newFakeRepo())
	id := uuid.New()
	_, err := uc.Update(context.Background(), service.UpdateCommand{
		ID:          id,
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      uuid.New(),
		StartDate:   "07-2025",
		EndDate:     strPtr("06-2025"),
	})
	require.ErrorIs(t, err, domain.ErrInvalidPeriod)
}

func strPtr(v string) *string { return &v }
