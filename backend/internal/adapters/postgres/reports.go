package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type TransactionCategory string

const (
	TransactionCategoryBeverages TransactionCategory = "beverages"
	TransactionCategorySnacks    TransactionCategory = "snacks"
	TransactionCategoryTobacco   TransactionCategory = "tobacco"
	TransactionCategoryGrocery   TransactionCategory = "grocery"
	TransactionCategoryLottery   TransactionCategory = "lottery"
	TransactionCategoryHousehold TransactionCategory = "household"
)

type TransactionStatus string

const (
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusRefunded  TransactionStatus = "refunded"
	TransactionStatusVoided    TransactionStatus = "voided"
)

type CategoryRevenueOrder string

const (
	CategoryRevenueOrderRevenue CategoryRevenueOrder = "revenue"
	CategoryRevenueOrderCount   CategoryRevenueOrder = "count"
)

type CategoryRevenueDirection string

const (
	CategoryRevenueDirectionAsc  CategoryRevenueDirection = "asc"
	CategoryRevenueDirectionDesc CategoryRevenueDirection = "desc"
)

type TransactionSummary struct {
	RevenueCents     int64 `db:"revenue_cents"`
	TransactionCount int64 `db:"transaction_count"`
	ItemsSold        int64 `db:"items_sold"`
	RefundedCount    int64 `db:"refunded_count"`
	VoidedCount      int64 `db:"voided_count"`
}

type DailyRevenue struct {
	Day          time.Time `db:"day"`
	RevenueCents int64     `db:"revenue_cents"`
}

type CategoryRevenue struct {
	Category         TransactionCategory `db:"category"`
	RevenueCents     int64               `db:"revenue_cents"`
	TransactionCount int64               `db:"transactions_count"`
	Total            int                 `db:"total"`
}

func (p postgres) GetTransactionSummary(ctx context.Context, storeID string, from, to time.Time) (TransactionSummary, error) {
	query := `
		SELECT
			COALESCE(SUM(t.total_cents) FILTER (WHERE t.status = $4), 0)::bigint AS revenue_cents,
			COUNT(t.id) FILTER (WHERE t.status = $4) AS transaction_count,
			COALESCE(SUM(t.items_count) FILTER (WHERE t.status = $4), 0) AS items_sold,
			COUNT(t.id) FILTER (WHERE t.status = $5) AS refunded_count,
			COUNT(t.id) FILTER (WHERE t.status = $6) AS voided_count
		FROM stores s
		LEFT JOIN transactions t
			ON t.store_id = s.id
			AND t.created_at >= ($2::date::timestamp AT TIME ZONE s.timezone)
			AND t.created_at < (($3::date + 1)::timestamp AT TIME ZONE s.timezone)
		WHERE s.id = $1
	`
	rows, err := p.pool.Query(
		ctx,
		query,
		storeID,
		civilDate(from),
		civilDate(to),
		TransactionStatusCompleted,
		TransactionStatusRefunded,
		TransactionStatusVoided,
	)
	if err != nil {
		return TransactionSummary{}, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[TransactionSummary])
}

func (p postgres) GetDailyRevenue(ctx context.Context, storeID string, from, to time.Time) ([]DailyRevenue, error) {
	query := `
		WITH bounds AS (
			SELECT
				s.id,
				s.timezone,
				($2::date::timestamp AT TIME ZONE s.timezone) AS from_ts,
				(($3::date + 1)::timestamp AT TIME ZONE s.timezone) AS to_ts
			FROM stores s
			WHERE s.id = $1
		),
		sales AS (
			SELECT
				(t.created_at AT TIME ZONE b.timezone)::date AS day,
				SUM(t.total_cents)::bigint AS revenue_cents
			FROM bounds b
			JOIN transactions t
				ON t.store_id = b.id
				AND t.status = $4
				AND t.created_at >= b.from_ts
				AND t.created_at < b.to_ts
			GROUP BY 1
		)
		SELECT
			series.day::date AS day,
			COALESCE(sales.revenue_cents, 0)::bigint AS revenue_cents
		FROM generate_series($2::date, $3::date, INTERVAL '1 day') AS series(day)
		LEFT JOIN sales ON sales.day = series.day::date
		ORDER BY series.day
	`
	rows, err := p.pool.Query(ctx, query, storeID, civilDate(from), civilDate(to), TransactionStatusCompleted)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[DailyRevenue])
}

func (p postgres) ListCategoryRevenue(ctx context.Context, storeID string, from, to time.Time, order CategoryRevenueOrder, direction CategoryRevenueDirection, limit, offset int) ([]CategoryRevenue, error) {
	query := fmt.Sprintf(`
		WITH bounds AS (
			SELECT
				s.id,
				($2::date::timestamp AT TIME ZONE s.timezone) AS from_ts,
				(($3::date + 1)::timestamp AT TIME ZONE s.timezone) AS to_ts
			FROM stores s
			WHERE s.id = $1
		),
		grouped AS (
			SELECT
				t.category,
				SUM(t.total_cents)::bigint AS revenue_cents,
				COUNT(*) AS transactions_count
			FROM bounds b
			JOIN transactions t
				ON t.store_id = b.id
				AND t.status = $4
				AND t.created_at >= b.from_ts
				AND t.created_at < b.to_ts
			GROUP BY t.category
		)
		SELECT
			category,
			revenue_cents,
			transactions_count,
			COUNT(*) OVER ()::int AS total
		FROM grouped
		ORDER BY %s, category ASC
		LIMIT $5 OFFSET $6
	`, categoryRevenueOrderBy(order, direction))

	rows, err := p.pool.Query(ctx, query, storeID, civilDate(from), civilDate(to), TransactionStatusCompleted, limit, offset)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[CategoryRevenue])
}

func (p postgres) CountRevenueCategories(ctx context.Context, storeID string, from, to time.Time) (int, error) {
	query := `
		WITH bounds AS (
			SELECT
				s.id,
				($2::date::timestamp AT TIME ZONE s.timezone) AS from_ts,
				(($3::date + 1)::timestamp AT TIME ZONE s.timezone) AS to_ts
			FROM stores s
			WHERE s.id = $1
		),
		grouped AS (
			SELECT t.category
			FROM bounds b
			JOIN transactions t
				ON t.store_id = b.id
				AND t.status = $4
				AND t.created_at >= b.from_ts
				AND t.created_at < b.to_ts
			GROUP BY t.category
		)
		SELECT COUNT(*)::int
		FROM grouped
	`
	var total int
	err := p.pool.QueryRow(ctx, query, storeID, civilDate(from), civilDate(to), TransactionStatusCompleted).
		Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func categoryRevenueOrderBy(order CategoryRevenueOrder, direction CategoryRevenueDirection) string {
	column := "revenue_cents"
	if order == CategoryRevenueOrderCount {
		column = "transactions_count"
	}

	sqlDirection := "DESC"
	if direction == CategoryRevenueDirectionAsc {
		sqlDirection = "ASC"
	}

	return column + " " + sqlDirection
}

func civilDate(t time.Time) string {
	return t.Format("2006-01-02")
}
