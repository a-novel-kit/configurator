package middlewares

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// UseContext replaces the default context of the request with a custom one. A timeout is also set to prevent the
// request from running indefinitely.
func UseContext(parentCTX context.Context, timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(parentCTX, timeout)
			defer cancel()

			go func() {
				<-ctx.Done()
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					w.WriteHeader(http.StatusGatewayTimeout)
				}
			}()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
