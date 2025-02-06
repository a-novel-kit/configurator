package middlewares_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/configurator/middlewares"
)

func TestUseContext(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "foo", "bar")

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value("foo") != "bar" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/something", nil)
		w := httptest.NewRecorder()

		middlewares.UseContext(ctx, time.Second)(handler).ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Timeout", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/something", nil)
		w := httptest.NewRecorder()

		middlewares.UseContext(context.Background(), time.Millisecond)(handler).ServeHTTP(w, req)
		require.Equal(t, http.StatusGatewayTimeout, w.Code)
	})
}
