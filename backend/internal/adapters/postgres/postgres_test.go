package postgres

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

const testDSN = "postgres://store:store@localhost:5432/store_sales_report_test?sslmode=disable"

type test struct {
	suite.Suite

	postgres postgres
}

func Test(t *testing.T) {
	suite.Run(t, new(test))
}

func (t *test) SetupSuite() {
	ctx := context.Background()

	admin, err := pgxpool.New(ctx, "postgres://store:store@localhost:5432/store_sales_report?sslmode=disable")
	t.Require().NoError(err)

	admin.Exec(ctx, `DROP DATABASE IF EXISTS store_sales_report_test`)
	admin.Exec(ctx, `CREATE DATABASE store_sales_report_test`)
	admin.Close()

	pool, err := pgxpool.New(ctx, testDSN)
	t.Require().NoError(err)

	t.applyMigrations()
	t.postgres = New(pool)
}

func (t *test) TearDownSuite() {
	admin, err := pgxpool.New(context.Background(), "postgres://store:store@localhost:5432/store_sales_report?sslmode=disable")
	t.Require().NoError(err)

	admin.Exec(context.Background(), `DROP DATABASE IF EXISTS store_sales_report_test`)
	admin.Close()

	t.postgres.pool.Close()
}

func (t *test) SetupTest() {
	_, err := t.postgres.pool.Exec(context.Background(), `TRUNCATE transactions, stores RESTART IDENTITY CASCADE`)
	t.Require().NoError(err)
}

func (t *test) applyMigrations() {
	_, file, _, ok := runtime.Caller(0)
	t.Require().True(ok)

	dir := filepath.Join(filepath.Dir(file), "..", "..", "..", "scripts", "migrations")
	m, err := migrate.New("file://"+dir, testDSN)
	t.Require().NoError(err)

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Require().NoError(err)
	}

	srcErr, dbErr := m.Close()
	t.Require().NoError(srcErr)
	t.Require().NoError(dbErr)
}

func (t *test) insertStore(id, name, currency, timezone string) {
	_, err := t.postgres.pool.Exec(context.Background(),
		`INSERT INTO stores (id, name, currency, timezone) VALUES ($1, $2, $3, $4)`,
		id, name, currency, timezone,
	)
	t.Require().NoError(err)
}

func (t *test) insertTx(storeID string, at time.Time, category string, items int, cents int64, status string) {
	_, err := t.postgres.pool.Exec(context.Background(), `
		INSERT INTO transactions (store_id, created_at, category, payment_method, items_count, total_cents, status)
		VALUES ($1, $2, $3, 'card', $4, $5, $6)
	`, storeID, at, category, items, cents, status)
	t.Require().NoError(err)
}

func (t *test) nyTime(day, hour, minute int) time.Time {
	return t.at("America/New_York", day, hour, minute)
}

func (t *test) johnsTime(day, hour, minute int) time.Time {
	return t.at("America/St_Johns", day, hour, minute)
}

func (t *test) at(timezone string, day, hour, minute int) time.Time {
	loc, err := time.LoadLocation(timezone)
	t.Require().NoError(err)
	return time.Date(2026, time.January, day, hour, minute, 0, 0, loc)
}
