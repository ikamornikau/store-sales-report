package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound = pgx.ErrNoRows
)

type Store struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	Currency string `db:"currency"`
	Timezone string `db:"timezone"`
}

func (p postgres) GetStores(ctx context.Context) ([]Store, error) {
	query := "SELECT id, name, currency, timezone FROM stores"
	rows, err := p.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[Store])
}

func (p postgres) GetStoreByID(ctx context.Context, id string) (Store, error) {
	query := "SELECT id, name, currency, timezone FROM stores WHERE id = $1"
	rows, err := p.pool.Query(ctx, query, id)
	if err != nil {
		return Store{}, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[Store])
}
