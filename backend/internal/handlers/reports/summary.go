package reports

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/ikamornikau/store-sales-report/backend/internal/adapters/postgres"
)

func (h reports) Summary(w http.ResponseWriter, r *http.Request) {
	from, to, err := parsePeriod(r.URL.Query())
	if err != nil {
		slog.WarnContext(r.Context(), "failed to parse period", slog.Any("error", err))
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

	summary, err := h.storage.GetTransactionSummary(r.Context(), store.ID, from, to)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to get transaction summary", slog.Any("error", err))
		http.Error(w, "failed to get transaction summary", http.StatusInternalServerError)
		return
	}

	writeJSON(w, r, summaryToAPI(store, summary))
}
