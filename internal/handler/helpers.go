package handler

import (
	"time"

	"github.com/seeques/test_junior/internal/models"
)

func parseMonthYear(s string) (time.Time, error) {
	return time.Parse("01-2006", s)
}

func toSubscriptionResponse(sub *models.Subscription) SubscriptionResponse {
	var endDate string
    if sub.EndDate != nil {
        endDate = sub.EndDate.Format("01-2006")
    }

	return SubscriptionResponse{
        ID:          sub.ID,
        ServiceName: sub.ServiceName,
        Price:       sub.Price,
        UserID:      sub.UserID.String(),
        StartDate:   sub.StartDate.Format("01-2006"),
        EndDate:     &endDate,
        CreatedAt:   sub.CreatedAt,
        UpdatedAt:   sub.UpdatedAt,
    }
}