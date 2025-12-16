package main

import (
	"context"

	"github.com/sapaude/go-shims/x/log"
	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化全局 Logger
	cfg := log.DefaultConfig()
	cfg.Level = logrus.DebugLevel
	cfg.Format = log.FormatJSON
	cfg.ReportCaller = true

	if err := log.InitGlobalLogger(cfg); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	// 演示 1：基本日志
	log.Infof("Application started")

	// 演示 2：带 Context 的日志（展示字段排序）
	ctx := log.WithRequestID(context.Background(), "req-12345")
	ctx = log.WithTracing(ctx, "trace-abc", "span-xyz")
	ctx = log.WithUserID(ctx, "user-789")

	// 演示 3：批量添加自定义字段
	ctx = log.WithFields(ctx, map[string]any{
		"order_id":   "order-456",
		"product_id": "prod-123",
		"amount":     99.99,
	})

	log.InfoContextf(ctx, "Order created successfully")

	// 演示 4：不同日志级别
	log.Debugf("Debug: detailed information")
	log.Warnf("Warning: something unusual happened")
	log.Errorf("Error: something went wrong")

	log.Infof("Application finished")
}
