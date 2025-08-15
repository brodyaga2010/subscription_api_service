package repository

import (
	"context"
	"eff_mobile/internal/model"
	"eff_mobile/pkg/pdb"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *model.Subscription) error
	Get(ctx context.Context, id int) (*model.Subscription, error)
	List(ctx context.Context) ([]model.Subscription, error)
	Update(ctx context.Context, sub *model.Subscription) error
	Delete(ctx context.Context, id int) error
	CalculateAmount(ctx context.Context, sub *model.SumRequest) (*model.SumResponse, error)
}

type subscriptionRepository struct {
	pool *pdb.Pool
}

func NewSubscriptionRepository(pool *pdb.Pool) SubscriptionRepository {
	return &subscriptionRepository{pool: pool}
}

func (r *subscriptionRepository) Create(ctx context.Context, subscription *model.Subscription) error {
	pool, err := r.pool.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("db pool get connection: %w", err)
	}

	query := `
        INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
        VALUES (@service_name, @price, @user_id, @start_date, @end_date)
        RETURNING id
    `
	args := pgx.NamedArgs{
		"service_name": subscription.Service,
		"price":        subscription.Price,
		"user_id":      subscription.UserID,
		"start_date":   subscription.StartDate,
		"end_date":     subscription.EndDate,
	}

	err = pool.QueryRow(ctx, query, args).Scan(&subscription.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrCreateSubscription
		}
		return fmt.Errorf("db query row: %w", err)
	}
	return nil
}

func (r *subscriptionRepository) Get(ctx context.Context, id int) (*model.Subscription, error) {
	pool, err := r.pool.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("db pool get connection: %w", err)
	}

	query := `SELECT service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = @id`
	subscription := &model.Subscription{}
	err = pool.QueryRow(ctx, query, pgx.NamedArgs{"id": id}).Scan(
		&subscription.Service,
		&subscription.Price,
		&subscription.UserID,
		&subscription.StartDate,
		&subscription.EndDate)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, model.ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("failed scan row: %v", err)
	}

	subscription.ID = id

	return subscription, nil
}

func (r *subscriptionRepository) Update(ctx context.Context, subscription *model.Subscription) error {
	pool, err := r.pool.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("db pool get connection: %w", err)
	}

	query := `
    UPDATE subscriptions 
    SET service_name = @name, price = @price, user_id = @user_id, start_date = @start, end_date = @stop 
    WHERE id = @id`
	args := pgx.NamedArgs{
		"name":    subscription.Service,
		"price":   subscription.Price,
		"user_id": subscription.UserID,
		"start":   subscription.StartDate,
		"stop":    subscription.EndDate,
		"id":      subscription.ID,
	}

	tag, err := pool.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed exec query: %v", err)
	}

	if tag.RowsAffected() == 0 {
		return model.ErrSubscriptionNotFound
	}

	return nil
}

func (r *subscriptionRepository) Delete(ctx context.Context, id int) error {
	pool, err := r.pool.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("db pool get connection: %w", err)
	}

	query := `DELETE FROM subscriptions WHERE id = @id`

	tag, err := pool.Exec(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return fmt.Errorf("db exec delete: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return model.ErrSubscriptionNotFound
	}
	return nil
}

func (r *subscriptionRepository) List(ctx context.Context) ([]model.Subscription, error) {
	pool, err := r.pool.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("db pool get connection: %w", err)
	}

	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("db query: %w", err)
	}
	defer rows.Close()

	var subscriptions []model.Subscription
	if err := pgxscan.ScanAll(&subscriptions, rows); err != nil {
		if err == pgx.ErrNoRows {
			return []model.Subscription{}, model.ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("failed to scan data: %v", err)
	}
	return subscriptions, nil
}

func (r *subscriptionRepository) CalculateAmount(ctx context.Context, req *model.SumRequest) (*model.SumResponse, error) {
	pool, err := r.pool.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("db pool get connection: %w", err)
	}

	var (
		conditions []string
		args       = pgx.NamedArgs{}
	)

	conditions = append(conditions, "start_date >= @from AND start_date <= @to")
	args["from"] = req.FromTime
	args["to"] = req.ToTime

	if req.UserID != "" {
		conditions = append(conditions, "user_id = @user_id")
		args["user_id"] = req.UserID
	}
	if req.ServiceName != "" {
		conditions = append(conditions, "service_name = @service_name")
		args["service_name"] = req.ServiceName
	}

	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE ` + strings.Join(conditions, " AND ")

	var total int

	err = pool.QueryRow(ctx, query, args).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("db query row: %w", err)
	}

	return &model.SumResponse{Total: total}, nil
}
