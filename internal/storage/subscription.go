package storage

import (
	"context"
	"fmt"

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
		return nil, fmt.Errorf("count subscriptions: %v", err)
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