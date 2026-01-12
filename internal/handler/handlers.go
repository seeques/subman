package handler

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/seeques/test_junior/internal/config"
	"github.com/seeques/test_junior/internal/storage"
)

type Handler struct {
	storage *storage.PostgresStorage
	cfg     config.Config
	logger  zerolog.Logger

}

func NewHandler(storage *storage.PostgresStorage, cfg config.Config, logger zerolog.Logger) *Handler {
	return &Handler{
		storage: storage,
		cfg:     cfg,
		logger:  logger.With().Str("component", "handler").Logger(),
	}
}

type SubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

type SubscriptionResponse struct {
	ID          string  `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
