package storage

import (
	"context"
	"time"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/google/uuid"
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

type TotalCostParams struct {
    StartPeriod time.Time
    EndPeriod   time.Time
    UserID      *uuid.UUID
    ServiceName string
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

func (s *PostgresStorage) GetSubscriptionsForPeriod(ctx context.Context, params TotalCostParams) ([]models.Subscription, error) {
	// start_date <= end_period ($2) and end_date >= start_period ($1)
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	FROM subscription
	WHERE start_date <= $2 AND (end_date >= $1 OR end_date IS NULL)`

	args := []interface{}{params.StartPeriod, params.EndPeriod}
	argNum := 3

	if params.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argNum)
		args = append(args, *params.UserID)
		argNum++
	}

	if params.ServiceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argNum)
		args = append(args, params.ServiceName)
		argNum++
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
        return nil, fmt.Errorf("get subscriptions for period: %w", err)
    }
    defer rows.Close()

	var subs []models.Subscription
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
            return nil, fmt.Errorf("scan subscription: %w", err)
        }
		subs = append(subs, sub)
	}
	return subs, rows.Err()
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
