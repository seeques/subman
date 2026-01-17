package handler

import (
	"encoding/json"
	"net/http"
	// "strings"
	"time"
	"log/slog"

	"github.com/google/uuid"
	"github.com/seeques/test_junior/internal/models"
	"github.com/seeques/test_junior/internal/response"
)

func parseMonthYear(s string) (time.Time, error) {
	return time.Parse("01-2006", s)
}

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
		response.RespondError(w, http.StatusInternalServerError, "failed to create log")
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
