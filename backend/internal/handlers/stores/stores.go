package stores

import (
	"context"
	"net/http"

	"github.com/ikamornikau/store-sales-report/backend/internal/adapters/postgres"
	"github.com/ikamornikau/store-sales-report/backend/internal/handlers"
	mw "github.com/ikamornikau/store-sales-report/backend/internal/middlewares"
)

type storage interface {
	GetStores(ctx context.Context) ([]postgres.Store, error)
}

type stores struct {
	storage storage
}

func NewHandler(storage storage) stores {
	return stores{storage: storage}
}

func (s stores) Mount(prefix string, mux *http.ServeMux, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle(handlers.MakeRoute(prefix, "GET", "/{$}"),
		mw.Chain(http.HandlerFunc(s.List), middlewares...),
	)
}
