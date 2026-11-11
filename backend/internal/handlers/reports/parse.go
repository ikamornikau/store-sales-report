package reports

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/ikamornikau/store-sales-report/backend/internal/adapters/postgres"
)

const (
	maxReportDays    = 92
	defaultPageLimit = 20
	maxPageLimit     = 100
)

func parsePeriod(query url.Values) (time.Time, time.Time, error) {
	rawFrom := query.Get("from")
	if rawFrom == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("param from is required")
	}

	rawTo := query.Get("to")
	if rawTo == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("param to is required")
	}

	from, err := time.Parse(time.DateOnly, rawFrom)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("param from should be a date in YYYY-MM-DD")
	}

	to, err := time.Parse(time.DateOnly, rawTo)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("param to should be a date in YYYY-MM-DD")
	}

	if from.After(to) {
		return time.Time{}, time.Time{}, fmt.Errorf("param from should not be after to")
	}

	if to.Sub(from).Hours() > float64(maxReportDays)*24 {
		return time.Time{}, time.Time{}, fmt.Errorf("date range should not exceed %d days", maxReportDays)
	}

	return from, to, nil
}

func parseSorting(query url.Values) (postgres.CategoryRevenueOrder, postgres.CategoryRevenueDirection, error) {
	order := postgres.CategoryRevenueOrder(query.Get("order_by"))
	if order == "" {
		order = postgres.CategoryRevenueOrderRevenue
	}

	switch order {
	case postgres.CategoryRevenueOrderRevenue, postgres.CategoryRevenueOrderCount:
	default:
		return "", "", fmt.Errorf("order should be revenue or count")
	}

	direction := postgres.CategoryRevenueDirection(query.Get("direction"))
	if direction == "" {
		direction = postgres.CategoryRevenueDirectionDesc
	}

	switch direction {
	case postgres.CategoryRevenueDirectionAsc, postgres.CategoryRevenueDirectionDesc:
	default:
		return "", "", fmt.Errorf("direction should be asc or desc")
	}

	return order, direction, nil
}

func parseLimitOffset(query url.Values) (int, int, error) {
	limit := defaultPageLimit
	if raw := query.Get("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 || parsed > maxPageLimit {
			return 0, 0, fmt.Errorf("limit should be an integer between 1 and %d", maxPageLimit)
		}

		limit = parsed
	}

	var offset int
	if raw := query.Get("offset"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 0 {
			return 0, 0, fmt.Errorf("offset should be a non-negative integer")
		}

		offset = parsed
	}

	return limit, offset, nil
}
