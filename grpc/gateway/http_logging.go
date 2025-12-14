package gateway_middleware

import (
	"net/http"
	"time"

	"github.com/sapaude/go-shims/x/log"
)

// HTTPLoggingMiddleware 记录 HTTP 请求的日志
func HTTPLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()
		// log.InfoContextf(ctx, "HTTP Request: Method=%s, URL=%s, RemoteAddr=%s, StartTime=%s",
		//     r.Method, r.URL.Path, r.RemoteAddr, start.Format(time.RFC3339))

		// 调用下一个处理器（可能是 gRPC-Gateway 的 ServeHTTP 方法，或者其他中间件）
		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.InfoContextf(ctx, "HTTP Response: Method=%s, URL=%s, Duration=%s",
			r.Method, r.URL.Path, duration)
	})
}
