package postgres

import (
	"context"
	"time"
)

func (t *test) TestGetTransactionSummary() {
	t.insertSales()

	summary, err := t.postgres.GetTransactionSummary(context.Background(), "S-1", civil(15), civil(16))
	t.Require().NoError(err)
	t.Require().Equal(TransactionSummary{
		RevenueCents:     4200,
		TransactionCount: 5,
		ItemsSold:        9,
		RefundedCount:    1,
		VoidedCount:      1,
	}, summary)
}

func (t *test) TestGetTransactionSummaryEmpty() {
	t.insertStore("S-1", "Manhattan", "USD", "America/New_York")

	summary, err := t.postgres.GetTransactionSummary(context.Background(), "S-1", civil(15), civil(16))
	t.Require().NoError(err)
	t.Require().Equal(TransactionSummary{}, summary)
}

func (t *test) TestGetDailyRevenue() {
	t.insertSales()

	days, err := t.postgres.GetDailyRevenue(context.Background(), "S-1", civil(13), civil(16))
	t.Require().NoError(err)
	t.Require().Equal([]DailyRevenue{
		{Day: civil(13), RevenueCents: 0},
		{Day: civil(14), RevenueCents: 100},
		{Day: civil(15), RevenueCents: 3000},
		{Day: civil(16), RevenueCents: 1200},
	}, days)
}

func (t *test) TestListCategoryRevenue() {
	t.insertSales()
	ctx := context.Background()

	byRevenue, err := t.postgres.ListCategoryRevenue(ctx, "S-1", civil(15), civil(16), CategoryRevenueOrderRevenue, CategoryRevenueDirectionDesc, 10, 0)
	t.Require().NoError(err)
	t.Require().Equal([]CategoryRevenue{
		{Category: TransactionCategoryGrocery, RevenueCents: 2000, TransactionCount: 1, Total: 4},
		{Category: TransactionCategoryBeverages, RevenueCents: 1500, TransactionCount: 2, Total: 4},
		{Category: TransactionCategoryHousehold, RevenueCents: 400, TransactionCount: 1, Total: 4},
		{Category: TransactionCategoryLottery, RevenueCents: 300, TransactionCount: 1, Total: 4},
	}, byRevenue)

	byCount, err := t.postgres.ListCategoryRevenue(ctx, "S-1", civil(15), civil(16), CategoryRevenueOrderCount, CategoryRevenueDirectionAsc, 10, 0)
	t.Require().NoError(err)
	t.Require().Equal([]CategoryRevenue{
		{Category: TransactionCategoryGrocery, RevenueCents: 2000, TransactionCount: 1, Total: 4},
		{Category: TransactionCategoryHousehold, RevenueCents: 400, TransactionCount: 1, Total: 4},
		{Category: TransactionCategoryLottery, RevenueCents: 300, TransactionCount: 1, Total: 4},
		{Category: TransactionCategoryBeverages, RevenueCents: 1500, TransactionCount: 2, Total: 4},
	}, byCount)

	page, err := t.postgres.ListCategoryRevenue(ctx, "S-1", civil(15), civil(16), CategoryRevenueOrderRevenue, CategoryRevenueDirectionDesc, 2, 1)
	t.Require().NoError(err)
	t.Require().Equal([]CategoryRevenue{
		{Category: TransactionCategoryBeverages, RevenueCents: 1500, TransactionCount: 2, Total: 4},
		{Category: TransactionCategoryHousehold, RevenueCents: 400, TransactionCount: 1, Total: 4},
	}, page)
}

func (t *test) TestCanadianSales() {
	t.insertSales()
	ctx := context.Background()

	store, err := t.postgres.GetStoreByID(ctx, "S-3")
	t.Require().NoError(err)
	t.Require().Equal(Store{
		ID: "S-3", Name: "Harbour", Currency: "CAD", Timezone: "America/St_Johns",
	}, store)

	summary, err := t.postgres.GetTransactionSummary(ctx, "S-3", civil(15), civil(16))
	t.Require().NoError(err)
	t.Require().Equal(TransactionSummary{
		RevenueCents:     3300,
		TransactionCount: 2,
		ItemsSold:        3,
		RefundedCount:    1,
	}, summary)

	days, err := t.postgres.GetDailyRevenue(ctx, "S-3", civil(14), civil(16))
	t.Require().NoError(err)
	t.Require().Equal([]DailyRevenue{
		{Day: civil(14), RevenueCents: 150},
		{Day: civil(15), RevenueCents: 1100},
		{Day: civil(16), RevenueCents: 2200},
	}, days)
}

func (t *test) TestCountRevenueCategories() {
	t.insertSales()

	total, err := t.postgres.CountRevenueCategories(context.Background(), "S-1", civil(15), civil(16))
	t.Require().NoError(err)
	t.Require().Equal(4, total)
}

func (t *test) TestReportsForUnknownStore() {
	ctx := context.Background()

	summary, err := t.postgres.GetTransactionSummary(ctx, "S-404", civil(15), civil(16))
	t.Require().NoError(err)
	t.Require().Equal(TransactionSummary{}, summary)

	days, err := t.postgres.GetDailyRevenue(ctx, "S-404", civil(15), civil(16))
	t.Require().NoError(err)
	t.Require().Equal([]DailyRevenue{
		{Day: civil(15), RevenueCents: 0},
		{Day: civil(16), RevenueCents: 0},
	}, days)

	categories, err := t.postgres.ListCategoryRevenue(ctx, "S-404", civil(15), civil(16), CategoryRevenueOrderRevenue, CategoryRevenueDirectionDesc, 10, 0)
	t.Require().NoError(err)
	t.Require().Empty(categories)

	total, err := t.postgres.CountRevenueCategories(ctx, "S-404", civil(15), civil(16))
	t.Require().NoError(err)
	t.Require().Equal(0, total)
}

func (t *test) insertSales() {
	t.insertStore("S-1", "Manhattan", "USD", "America/New_York")
	t.insertStore("S-2", "Other", "USD", "America/Chicago")
	t.insertStore("S-3", "Harbour", "CAD", "America/St_Johns")

	t.insertTx("S-1", t.nyTime(14, 23, 59), "beverages", 1, 100, "completed")
	t.insertTx("S-1", t.nyTime(15, 0, 0), "beverages", 2, 1000, "completed")
	t.insertTx("S-1", t.nyTime(15, 12, 0), "snacks", 4, 800, "refunded")
	t.insertTx("S-1", t.nyTime(15, 13, 0), "tobacco", 5, 700, "voided")
	t.insertTx("S-1", t.nyTime(15, 23, 59), "grocery", 3, 2000, "completed")
	t.insertTx("S-1", t.nyTime(16, 0, 0), "lottery", 1, 300, "completed")
	t.insertTx("S-1", t.nyTime(16, 10, 0), "beverages", 1, 500, "completed")
	t.insertTx("S-1", t.nyTime(16, 18, 0), "household", 2, 400, "completed")
	t.insertTx("S-1", t.nyTime(17, 0, 0), "household", 1, 9000, "completed")
	t.insertTx("S-2", t.nyTime(15, 12, 0), "beverages", 10, 5000, "completed")

	t.insertTx("S-3", t.johnsTime(14, 23, 59), "beverages", 1, 150, "completed")
	t.insertTx("S-3", t.johnsTime(15, 0, 0), "beverages", 2, 1100, "completed")
	t.insertTx("S-3", t.johnsTime(15, 12, 0), "snacks", 1, 250, "refunded")
	t.insertTx("S-3", t.johnsTime(16, 0, 0), "grocery", 1, 2200, "completed")
	t.insertTx("S-3", t.johnsTime(17, 0, 0), "lottery", 1, 8000, "completed")
}

func civil(day int) time.Time {
	return time.Date(2026, time.January, day, 0, 0, 0, 0, time.UTC)
}
