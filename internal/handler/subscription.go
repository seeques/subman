package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"errors"
	
	"github.com/google/uuid"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/seeques/test_junior/internal/models"
	"github.com/seeques/test_junior/internal/response"
	"github.com/seeques/test_junior/internal/storage"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req SubscriptionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Warn("invalid JSON in request body", "error", err)
		response.RespondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// Validation
	if req.Price <= 0 {
		response.RespondError(w, http.StatusBadRequest, "price must be more than zero")
		return
	}

	if req.ServiceName == "" || req.UserID == "" || req.StartDate == "" {
		response.RespondError(w, http.StatusBadRequest, "request field is empty")
		return
	}

	// Parse uuid
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		response.RespondError(w, http.StatusBadRequest, "invalid user_id, must be UUID")
		return
	}

	// Parse dates to check if they match MM-YYYY format
	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		response.RespondError(w, http.StatusBadRequest, "invalid start_date, expected MM-YYYY")
		return
	}

	var endDate *time.Time
	if req.EndDate != "" {
		parsed, err := parseMonthYear(req.EndDate)
		if err != nil {
			response.RespondError(w, http.StatusBadRequest, "invalid end_date, expected MM-YYYY")
			return
		}
		endDate = &parsed
	}

	ctx := r.Context()
	sub := &models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	// Create new subscription
	err = h.storage.CreateSubscription(ctx, sub)
	if err != nil {
		slog.Error("failed to create subscription",
			"error", err,
			"service_name", req.ServiceName,
		)
		response.RespondError(w, http.StatusInternalServerError, "failed to create subscriptions")
		return
	}

	slog.Info("subscription created",
		"subscription_id", sub.ID,
		"service_name", sub.ServiceName,
	)

	var subEndDate string
	if sub.EndDate != nil {
		subEndDate = sub.EndDate.Format("01-2006")
	}

	// Make a response
	response.RespondJSON(w, http.StatusCreated, SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID.String(),
		StartDate:   sub.StartDate.Format("01-2006"),
		EndDate:     &subEndDate,
		CreatedAt:   sub.CreatedAt,
		UpdatedAt:   sub.UpdatedAt,
	})
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.RespondError(w, http.StatusBadRequest, "invalid id")
        return
	}

	sub, err := h.storage.GetSubscription(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.RespondError(w, http.StatusNotFound, "subbscription not found")
			return
		}
		slog.Error("failed to get subscription", "error", err, "id", id)
		response.RespondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.RespondJSON(w, http.StatusOK, toSubscriptionResponse(sub))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	// default to 1 page
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	// default limit to 10
	if err != nil || limit < 1 {
		limit = 10
	}

	// cap limit to 100
	if limit > 100 {
		limit = 100
	}

	result, err := h.storage.ListAllSubscriptions(r.Context(), storage.ListParams{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		slog.Error("failed to list subscriptions",
			"error", err)
		response.RespondError(w, http.StatusInternalServerError, "failed to list subscriptions")
		return
	}

	totalPages := (result.Total + limit - 1) / limit

	response.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"data": result.Subscriptions,
		"meta": map[string]int{
			"page":        page,
			"limit":       limit,
			"total":       result.Total,
			"total_pages": totalPages,
		},
	})
}
