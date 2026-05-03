package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Romasmi/social-network/internal/utils"
	"github.com/gorilla/mux"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var rw *utils.StatusResponseWriter
		if existing, ok := w.(*utils.StatusResponseWriter); ok {
			rw = existing
		} else {
			rw = utils.NewStatusResponseWriter(w)
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		path := r.URL.Path
		if route := mux.CurrentRoute(r); route != nil {
			if tpl, err := route.GetPathTemplate(); err == nil && tpl != "" {
				path = tpl
			}
		}

		slog.Info("request processed",
			slog.String("method", r.Method),
			slog.String("path", path),
			slog.Int("status", rw.StatusCode),
			slog.Duration("duration", duration),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		)
	})
}
