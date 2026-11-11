package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type postgres struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) postgres {
	return postgres{pool: pool}
}
