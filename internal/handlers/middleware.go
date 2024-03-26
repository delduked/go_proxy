package handlers

import (
	"net/http"

	l "prx/internal/logger"
)

type wrapperWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrapperWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

// validateUpdateRequest ensures that the request is a POST and Content-Type is application/json.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Log.Info(r.Method, r.URL.Path)
		wrapped := &wrapperWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		l.Log.Info("New Request:","method",r.Method, "status",wrapped.statusCode,"path", r.URL.Path,"host", r.Host)
		next.ServeHTTP(wrapped, r)
	})
}

type Middleware func(http.Handler) http.Handler

func MiddlewareStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}
