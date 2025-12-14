package gateway_middleware

import (
	"net/http"

	"github.com/sapaude/go-shims/x/log"
)

// HTTPAuthMiddleware 简单的 HTTP 认证中间件
func HTTPAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求头中的 Authorization
		token := r.Header.Get("Authorization")
		if token == "" || token != "Bearer my-http-secret-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Infof("HTTP Auth Middleware: URL=%s - Authenticated (for demo purposes)", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
