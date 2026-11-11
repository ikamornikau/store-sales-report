package reports

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/ikamornikau/store-sales-report/backend/internal/adapters/postgres"
	"github.com/ikamornikau/store-sales-report/backend/internal/handlers"
	mw "github.com/ikamornikau/store-sales-report/backend/internal/middlewares"
)

type storage interface {
	GetStoreByID(ctx context.Context, id string) (postgres.Store, error)
	GetTransactionSummary(ctx context.Context, storeID string, from, to time.Time) (postgres.TransactionSummary, error)
	GetDailyRevenue(ctx context.Context, storeID string, from, to time.Time) ([]postgres.DailyRevenue, error)
	ListCategoryRevenue(ctx context.Context, storeID string, from, to time.Time, order postgres.CategoryRevenueOrder, direction postgres.CategoryRevenueDirection, limit, offset int) ([]postgres.CategoryRevenue, error)
	CountRevenueCategories(ctx context.Context, storeID string, from, to time.Time) (int, error)
}

type reports struct {
	storage storage
}

func NewHandler(storage storage) reports {
	return reports{storage: storage}
}

func (h reports) Mount(prefix string, mux *http.ServeMux, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle(handlers.MakeRoute(prefix, "GET", "/{store_id}/summary"),
		mw.Chain(http.HandlerFunc(h.Summary), middlewares...),
	)
	mux.Handle(handlers.MakeRoute(prefix, "GET", "/{store_id}/revenue/daily"),
		mw.Chain(http.HandlerFunc(h.DailyRevenue), middlewares...),
	)
	mux.Handle(handlers.MakeRoute(prefix, "GET", "/{store_id}/revenue/category"),
		mw.Chain(http.HandlerFunc(h.GetCategoryRevenue), middlewares...),
	)
}

func writeJSON(w http.ResponseWriter, r *http.Request, payload any) {
	res, err := json.Marshal(payload)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to marshal response", slog.Any("error", err))
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
