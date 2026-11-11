package stores

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (s stores) List(w http.ResponseWriter, r *http.Request) {
	stores, err := s.storage.GetStores(r.Context())
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to get stores", slog.Any("error", err))
		http.Error(w, "failed to get stores", http.StatusInternalServerError)
		return
	}

	apiStores := storesToAPI(stores)
	res, err := json.Marshal(apiStores)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to marshal stores", slog.Any("error", err))
		http.Error(w, "failed to marshal stores", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
