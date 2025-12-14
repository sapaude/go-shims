package gateway_middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sapaude/go-shims/x/log"
)

// DynamicContextMiddleware 动态注入值到ctx
func DynamicContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCtx := r.Context()
		ctx := log.WithTraceID(reqCtx, uuid.NewString())

		ctx = log.WithUserID(ctx, "uid-100")
		ctx = log.WithCustomField(ctx, "x-host", "localhost")
		ctx = log.WithCustomField(ctx, "x-version", "v1.0.0")
		ctx = log.WithCustomField(ctx, "x-uid", 100)
		newReq := r.WithContext(ctx)

		// log.InfoContextf(ctx, "middleware log... request ctx: %v", c.Request().Context())

		// 继续调用下游处理器
		next.ServeHTTP(w, newReq)
	})
}
