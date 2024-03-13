package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func recordResponse(method, target string) *httptest.ResponseRecorder {
	mux := chi.NewMux()
	NewHealthHandler(mux)
	req := httptest.NewRequest(method, target, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestHealthHandler(t *testing.T) {
	t.Run("Health endpoint success case", func(t *testing.T) {
		rec := recordResponse(http.MethodGet, "/health")
		res := rec.Result()

		assert.Equal(t, res.StatusCode, http.StatusOK)
	})

	t.Run("Readiness endpoint success case", func(t *testing.T) {
		rec := recordResponse(http.MethodGet, "/readiness")
		res := rec.Result()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
