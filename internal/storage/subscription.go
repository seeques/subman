package storage

import (
	"fmt"
	"context"
	
	"github.com/seeques/test_junior/internal/models"
)

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