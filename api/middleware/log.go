package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// func Logging(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		start := time.Now()
// 		next.ServeHTTP(w, req)
// 		log.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
// 	})
// }

func Logging(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Infow("New request",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
			"time", time.Since(start),
		)
	})
}
