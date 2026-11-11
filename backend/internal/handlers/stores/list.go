package stores

import "net/http"

func (s stores) List(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from stores list handler!"))
}
