package reports

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/ikamornikau/store-sales-report/backend/internal/adapters/postgres"
)

func (h reports) GetCategoryRevenue(w http.ResponseWriter, r *http.Request) {
	from, to, err := parsePeriod(r.URL.Query())
	if err != nil {
		slog.WarnContext(r.Context(), "failed to parse period", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, direction, err := parseSorting(r.URL.Query())
	if err != nil {
		slog.WarnContext(r.Context(), "failed to parse sorting", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limit, offset, err := parseLimitOffset(r.URL.Query())
	if err != nil {
		slog.WarnContext(r.Context(), "failed to parse limit offset", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	store, err := h.storage.GetStoreByID(r.Context(), r.PathValue("store_id"))
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			slog.WarnContext(r.Context(), "store not found", slog.Any("error", err))
			http.Error(w, "store not found", http.StatusNotFound)
			return
		}

		slog.ErrorContext(r.Context(), "failed to get store", slog.Any("error", err))
		http.Error(w, "failed to get store", http.StatusInternalServerError)
		return
	}

	rows, err := h.storage.ListCategoryRevenue(r.Context(), store.ID, from, to, order, direction, limit, offset)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to list category revenue", slog.Any("error", err))
		http.Error(w, "failed to list category revenue", http.StatusInternalServerError)
		return
	}

	var total int
	if len(rows) < 1 {
		total, err = h.storage.CountRevenueCategories(r.Context(), store.ID, from, to)
		if err != nil {
			slog.ErrorContext(r.Context(), "failed to count revenue categories", slog.Any("error", err))
			http.Error(w, "failed to count revenue categories", http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, r, categoriesToAPI(store, rows, total))
}
