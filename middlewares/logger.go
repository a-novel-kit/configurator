package middlewares

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

// ZeroLogger sets a zerolog logger to the request context and logs the request when it's done.
//
// It is recommended to set middleware.RequestID before this middleware to have a unique identifier for each request.
func ZeroLogger(parentLogger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Capture the response writer to get the status code.
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			ctx := r.Context()

			requestLogger := parentLogger.With().
				// Used to group requests together.
				Str("request_id", middleware.GetReqID(ctx)).
				Logger()

			// Inject the logger into the following middlewares context.
			next.ServeHTTP(ww, r.WithContext(requestLogger.WithContext(ctx)))

			// Switch severity based on status code.
			status := ww.Status()
			event := lo.
				If(status > 499, requestLogger.Error).
				ElseIf(status > 399, requestLogger.Warn).
				Else(requestLogger.Info)

			// Write the actual log.
			event().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Int("status", status).
				Str("user_agent", r.UserAgent()).
				Str("remote", r.RemoteAddr).
				Str("path", r.URL.Path).
				Msg("request completed")
		})
	}
}
