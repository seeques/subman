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
    ServiceName string `json:"service_name" example:"Yandex Plus"`
    Price       int    `json:"price" example:"400"`
    UserID      string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
    StartDate   string `json:"start_date" example:"07-2025"`
    EndDate     string `json:"end_date,omitempty" example:"12-2025"`
}

type SubscriptionResponse struct {
    ID          int       `json:"id" example:"1"`
    ServiceName string    `json:"service_name" example:"Yandex Plus"`
    Price       int       `json:"price" example:"400"`
    UserID      string    `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
    StartDate   string    `json:"start_date" example:"07-2025"`
    EndDate     *string   `json:"end_date,omitempty" example:"12-2025"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type ErrorResponse struct {
    Error string `json:"error" example:"invalid request"`
}

type ListResponse struct {
    Data []SubscriptionResponse `json:"data"`
    Meta ListMeta               `json:"meta"`
}

type ListMeta struct {
    Page       int `json:"page" example:"1"`
    Limit      int `json:"limit" example:"10"`
    Total      int `json:"total" example:"100"`
    TotalPages int `json:"total_pages" example:"10"`
}

type TotalCostResponse struct {
    TotalCost          int    `json:"total_cost" example:"3600"`
    Currency           string `json:"currency" example:"RUB"`
    PeriodStart        string `json:"period_start" example:"01-2025"`
    PeriodEnd          string `json:"period_end" example:"06-2025"`
    SubscriptionsCount int    `json:"subscriptions_count" example:"3"`
}