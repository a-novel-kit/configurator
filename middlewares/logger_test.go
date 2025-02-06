package middlewares_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/configurator/middlewares"
)

func TestZeroLogger(t *testing.T) {
	testCases := []struct {
		name string

		status int

		expect string
	}{
		{
			name:   "OK",
			status: http.StatusOK,
			expect: `{"level":"info","message":"request completed","method":"GET","path":"/something",
"remote":"127.0.0.1:1234","request_id":"","status":200,"url":"/something","user_agent":"test-ua"}`,
		},
		{
			name:   "NotFound",
			status: http.StatusNotFound,
			expect: `{"level":"warn","message":"request completed","method":"GET","path":"/something",
"remote":"127.0.0.1:1234","request_id":"","status":404,"url":"/something","user_agent":"test-ua"}`,
		},
		{
			name:   "InternalServerError",
			status: http.StatusInternalServerError,
			expect: `{"level":"error","message":"request completed","method":"GET","path":"/something",
"remote":"127.0.0.1:1234","request_id":"","status":500,"url":"/something","user_agent":"test-ua"}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var buff bytes.Buffer
			logger := zerolog.New(&buff)

			baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(testCase.status)
			})

			baseHandlerWithRequestID := middleware.RequestID(baseHandler)
			baseHandlerWithLogger := middlewares.ZeroLogger(&logger)(baseHandlerWithRequestID)

			req := httptest.NewRequest(http.MethodGet, "/something", nil)
			req.Header.Set("User-Agent", "test-ua")
			req.URL.Path = "/something"
			req.RequestURI = ""
			req.RemoteAddr = "127.0.0.1:1234"

			w := httptest.NewRecorder()
			baseHandlerWithLogger.ServeHTTP(w, req)

			require.JSONEq(t, testCase.expect, buff.String())
		})
	}
}
