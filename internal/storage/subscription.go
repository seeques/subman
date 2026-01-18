package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/seeques/test_junior/internal/models"
)

type ListParams struct {
	Page  int
	Limit int
}

type ListResult struct {
	Subscriptions []models.Subscription
	Total         int
}

func (s *PostgresStorage) CreateSubscription(ctx context.Context, sub *models.Subscription) error {
	query := `INSERT INTO subscription (service_name, price, user_id, start_date, end_date)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at`

	err := s.pool.QueryRow(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create subscription: %w", err)
	}
	return nil
}

func (s *PostgresStorage) GetSubscription(ctx context.Context, id int) (*models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	FROM subscription
	WHERE id = $1`

	var sub models.Subscription
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}
	return &sub, nil
}

func (s *PostgresStorage) UpdateSubscription(ctx context.Context, sub *models.Subscription) error {
	query := `UPDATE subscription SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5, updated_at = NOW()
	WHERE id = $6
	RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at`

	err := s.pool.QueryRow(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.ID).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		fmt.Errorf("update subscription: %w", err)
	}
	return nil
}

func (s *PostgresStorage) DeleteSubscription(ctx context.Context, id int)  error {
	query := `DELETE FROM subscription WHERE id = $1`

	result, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete subscription", err)
	}

	if result.RowsAffected() == 0 {
        return pgx.ErrNoRows
    }

	return nil
}

func (s *PostgresStorage) ListAllSubscriptions(ctx context.Context, params ListParams) (*ListResult, error) {
	// limit = 10
	// 1st page: offset = 0
	// 2nd page: offset = 10
	offset := (params.Page - 1) * params.Limit

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM subscription`
	err := s.pool.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count subscriptions: %w", err)
	}

	// Get page
	pageQuery := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	FROM subscription
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2`

	rows, err := s.pool.Query(ctx, pageQuery, params.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list subscription: %v", err)
	}
	defer rows.Close()

	var subscriptions []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan subscription: %v", err)
		}
		subscriptions = append(subscriptions, sub)
	}

	return &ListResult{
		Subscriptions: subscriptions,
		Total:         total,
	}, nil
}
