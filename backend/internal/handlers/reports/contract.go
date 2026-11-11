package reports

import (
	"math"
	"time"

	"github.com/ikamornikau/store-sales-report/backend/internal/adapters/postgres"
)

type Summary struct {
	Currency           string `json:"currency"`
	RevenueCents       int64  `json:"revenue_cents"`
	TransactionCount   int64  `json:"transaction_count"`
	RefundedCount      int64  `json:"refunded_count"`
	VoidedCount        int64  `json:"voided_count"`
	AverageBasketCents int64  `json:"average_basket_cents,omitempty"`
	ItemsSold          int64  `json:"items_sold"`
}

type DailyRevenue struct {
	Currency string       `json:"currency"`
	Timezone string       `json:"timezone"`
	Points   []DailyPoint `json:"points"`
}

type DailyPoint struct {
	Day          string `json:"day"`
	RevenueCents int64  `json:"revenue_cents"`
}

type CategoryRevenuePage struct {
	Currency string            `json:"currency"`
	Total    int               `json:"total"`
	Items    []CategoryRevenue `json:"items"`
}

type CategoryRevenue struct {
	Category         string `json:"category"`
	RevenueCents     int64  `json:"revenue_cents"`
	TransactionCount int64  `json:"transaction_count"`
}

func summaryToAPI(store postgres.Store, summary postgres.TransactionSummary) Summary {
	return Summary{
		Currency:           store.Currency,
		RevenueCents:       summary.RevenueCents,
		TransactionCount:   summary.TransactionCount,
		RefundedCount:      summary.RefundedCount,
		VoidedCount:        summary.VoidedCount,
		AverageBasketCents: averageBasketCents(summary.RevenueCents, summary.TransactionCount),
		ItemsSold:          summary.ItemsSold,
	}
}

func averageBasketCents(revenueCents, transactionCount int64) int64 {
	if transactionCount == 0 {
		return 0
	}

	return int64(math.Round(float64(revenueCents) / float64(transactionCount)))
}

func dailyToAPI(store postgres.Store, rows []postgres.DailyRevenue) DailyRevenue {
	points := make([]DailyPoint, 0, len(rows))
	for _, row := range rows {
		points = append(points, DailyPoint{
			Day:          row.Day.Format(time.DateOnly),
			RevenueCents: row.RevenueCents,
		})
	}

	return DailyRevenue{
		Currency: store.Currency,
		Timezone: store.Timezone,
		Points:   points,
	}
}

func categoriesToAPI(store postgres.Store, rows []postgres.CategoryRevenue, total int) CategoryRevenuePage {
	if len(rows) > 0 {
		total = rows[0].Total
	}

	items := make([]CategoryRevenue, 0, len(rows))
	for _, row := range rows {
		items = append(items, CategoryRevenue{
			Category:         string(row.Category),
			RevenueCents:     row.RevenueCents,
			TransactionCount: row.TransactionCount,
		})
	}

	return CategoryRevenuePage{
		Currency: store.Currency,
		Total:    total,
		Items:    items,
	}
}
