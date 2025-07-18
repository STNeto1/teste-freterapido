package http

import (
	"net/http"
)

// HealthHandler returns 200 OK
func (r *Router) HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		_ = err
	}
}
