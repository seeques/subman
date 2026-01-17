package handler

import (
	"time"

	"github.com/seeques/test_junior/internal/config"
	"github.com/seeques/test_junior/internal/storage"
)

type Handler struct {
	storage *storage.PostgresStorage
	cfg     config.Config
}

func NewHandler(storage *storage.PostgresStorage, cfg config.Config) *Handler {
	return &Handler{
		storage: storage,
		cfg:     cfg,
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
	ID          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      string    `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
