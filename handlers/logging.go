package handlers

import (
	"net/http"

	"github.com/rs/zerolog"
)

func LoggingMiddleware(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info().
				Str("proto", r.Proto).
				Str("uri", r.RequestURI).
				Str("method", r.Method).
				Str("remote", r.RemoteAddr).
				Str("user-agent", r.UserAgent()).
				Msg("request started")
			next.ServeHTTP(w, r)
		})
	}
}
