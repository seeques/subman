package handler

import (
	"time"

	"github.com/seeques/subman/internal/models"
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

func countMonths(start, end time.Time) int {
    years := end.Year() - start.Year()
    months := int(end.Month()) - int(start.Month())
    return years*12 + months + 1  // +1 because inclusive
}

func calculateTotalCost(subs []models.Subscription, startPeriod, endPeriod time.Time) int {
    total := 0

    for _, sub := range subs {
        // Find if startPeriod overlaps with the startDate
        overlapStart := sub.StartDate
        if startPeriod.After(overlapStart) {
            overlapStart = startPeriod
        }

        // Find if endDate overlaps with the endPeriod
        overlapEnd := endPeriod
        if sub.EndDate != nil && sub.EndDate.Before(overlapEnd) {
            overlapEnd = *sub.EndDate
        }

        if overlapStart.After(overlapEnd) {
            continue
        }

        months := countMonths(overlapStart, overlapEnd)
        total += months * sub.Price
    }

    return total
}