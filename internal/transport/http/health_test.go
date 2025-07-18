package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	// TODO: nil, nil, nil <- fix for a proper health check
	NewRouter(nil, nil, nil).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, status)
	}

	expected := "OK"
	if body := rr.Body.String(); body != expected {
		t.Errorf("expected body %q; got %q", expected, body)
	}
}
