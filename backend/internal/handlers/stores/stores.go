package stores

import (
	"fmt"
	"net/http"

	mw "github.com/ikamornikau/store-sales-report/backend/internal/middlewares"
)

type stores struct{}

func NewHandler() stores {
	return stores{}
}

func (s stores) Mount(prefix string, mux *http.ServeMux, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle(makeRoute(prefix, "GET", "/{$}"),
		mw.Chain(http.HandlerFunc(s.List), middlewares...),
	)
}

func makeRoute(prefix, method, pattern string) string {
	return fmt.Sprintf("%s %s%s", method, prefix, pattern)
}
