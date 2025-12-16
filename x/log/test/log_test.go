package test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/sapaude/go-shims/x/log"
	"github.com/sirupsen/logrus"
)

// TestGlobalLoggerInit 测试全局 Logger 初始化
func TestGlobalLoggerInit(t *testing.T) {
	cfg := log.DefaultConfig()
	cfg.Level = logrus.DebugLevel
	cfg.Format = log.FormatJSON

	err := log.InitGlobalLogger(cfg)
	if err != nil {
		t.Fatalf("InitGlobalLogger failed: %v", err)
	}

	if !log.IsGlobalLoggerInitialized() {
		t.Fatal("Global logger should be initialized")
	}

	// 测试基本日志方法
	log.Debugf("Debug message")
	log.Infof("Info message")
	log.Warnf("Warning message")
	log.Errorf("Error message")
}

// TestContextFields 测试 Context 字段提取和注入
func TestContextFields(t *testing.T) {
	buf := &bytes.Buffer{}
	cfg := log.DefaultConfig()
	cfg.Format = log.FormatJSON
	cfg.Output = buf
	cfg.ReportCaller = false // 关闭 caller 以简化输出

	logger, err := log.NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	defer logger.Close()

	// 使用 WithFields 批量添加字段
	ctx := log.WithRequestID(context.Background(), "req-12345")
	ctx = log.WithTracing(ctx, "trace-xyz", "span-abc")
	ctx = log.WithUserID(ctx, "user-456")
	ctx = log.WithFields(ctx, map[string]any{
		"custom_key1": "value1",
		"custom_key2": 123,
	})

	logger.InfoContextf(ctx, "Test message with context")

	output := buf.String()
	t.Logf("Output: %s", output)

	// 验证字段是否存在
	if !strings.Contains(output, `"request_id":"req-12345"`) {
		t.Error("request_id not found in output")
	}
	if !strings.Contains(output, `"trace_id":"trace-xyz"`) {
		t.Error("trace_id not found in output")
	}
	if !strings.Contains(output, `"span_id":"span-abc"`) {
		t.Error("span_id not found in output")
	}
	if !strings.Contains(output, `"user_id":"user-456"`) {
		t.Error("user_id not found in output")
	}
	if !strings.Contains(output, `"custom_key1":"value1"`) {
		t.Error("custom_key1 not found in output")
	}
}

// TestFieldOrdering 测试 JSON 字段顺序
func TestFieldOrdering(t *testing.T) {
	buf := &bytes.Buffer{}
	cfg := log.DefaultConfig()
	cfg.Format = log.FormatJSON
	cfg.Output = buf
	cfg.ReportCaller = true // 启用 caller 以测试 file 和 func 字段
	cfg.JSONPretty = false

	logger, err := log.NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	defer logger.Close()

	// 添加上下文字段
	ctx := log.WithRequestID(context.Background(), "req-123")
	ctx = log.WithCustomField(ctx, "custom", "data")

	logger.InfoContextf(ctx, "test message")

	output := buf.String()
	t.Logf("Output: %s", output)

	// 验证字段顺序：time 应该在 level 之前，level 在 file 之前，msg 应该在最后
	timeIdx := strings.Index(output, `"time"`)
	levelIdx := strings.Index(output, `"level"`)
	fileIdx := strings.Index(output, `"file"`)
	funcIdx := strings.Index(output, `"func"`)
	msgIdx := strings.Index(output, `"msg"`)

	if timeIdx == -1 || levelIdx == -1 || fileIdx == -1 || funcIdx == -1 || msgIdx == -1 {
		t.Fatal("Missing required fields in output")
	}

	if !(timeIdx < levelIdx && levelIdx < fileIdx && fileIdx < funcIdx) {
		t.Errorf("Field order incorrect: time(%d) < level(%d) < file(%d) < func(%d)",
			timeIdx, levelIdx, fileIdx, funcIdx)
	}

	// msg 应该在最后（在所有业务字段之后）
	if msgIdx < funcIdx {
		t.Errorf("msg should be last, but msg(%d) < func(%d)", msgIdx, funcIdx)
	}

	// 验证可以正常解析为 JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &jsonData); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// 验证必需字段存在
	requiredFields := []string{"time", "level", "file", "func", "msg"}
	for _, field := range requiredFields {
		if _, ok := jsonData[field]; !ok {
			t.Errorf("Required field '%s' not found in JSON", field)
		}
	}
}

// TestResourceCleanup 测试资源清理（文件句柄关闭）
func TestResourceCleanup(t *testing.T) {
	tmpFile := "/tmp/test_log_cleanup.log"
	defer os.Remove(tmpFile) // 确保清理临时文件

	cfg := log.DefaultConfig()
	cfg.FilePath = tmpFile
	cfg.Format = log.FormatJSON

	logger, err := log.NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	logger.Infof("test message")

	// 关闭 Logger
	err = logger.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// 验证文件可以被删除（说明句柄已关闭）
	err = os.Remove(tmpFile)
	if err != nil {
		t.Errorf("Failed to remove log file (handle may not be closed): %v", err)
	}
}

// TestDynamicConfig 测试动态配置修改
func TestDynamicConfig(t *testing.T) {
	buf := &bytes.Buffer{}
	cfg := log.DefaultConfig()
	cfg.Format = log.FormatJSON
	cfg.Output = buf
	cfg.Level = logrus.InfoLevel
	cfg.ReportCaller = false

	logger, err := log.NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	defer logger.Close()

	// Debug 消息不应该输出（level 是 Info）
	logger.Debugf("This should not appear")
	if buf.Len() > 0 {
		t.Error("Debug message should not appear when level is Info")
	}

	// 动态修改为 Debug 级别
	logger.SetLevel(logrus.DebugLevel)
	logger.Debugf("This should appear")
	if buf.Len() == 0 {
		t.Error("Debug message should appear after changing level to Debug")
	}

	// 测试格式切换
	buf.Reset()
	logger.SetFormatter(log.FormatText)
	logger.Infof("Text format message")

	output := buf.String()
	// Text 格式不应该包含 JSON 结构
	if strings.Contains(output, `"level"`) {
		t.Error("Text format should not contain JSON structure")
	}
}

// TestFileOutput 测试文件输出
func TestFileOutput(t *testing.T) {
	tmpFile := "/tmp/test_log_file.log"
	defer os.Remove(tmpFile)

	cfg := log.DefaultConfig()
	cfg.FilePath = tmpFile
	cfg.Format = log.FormatJSON
	cfg.ReportCaller = true

	logger, err := log.NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	ctx := log.WithRequestID(context.Background(), "req-file-test")
	logger.InfoContextf(ctx, "Message written to file")

	// 关闭 logger 以确保内容被刷新
	logger.Close()

	// 读取文件验证内容
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, `"msg":"Message written to file"`) {
		t.Error("Log message not found in file")
	}
	if !strings.Contains(output, `"request_id":"req-file-test"`) {
		t.Error("request_id not found in file")
	}
}

// TestWithFieldsPerformance 测试 WithFields 批量操作
func TestWithFieldsPerformance(t *testing.T) {
	buf := &bytes.Buffer{}
	cfg := log.DefaultConfig()
	cfg.Format = log.FormatJSON
	cfg.Output = buf
	cfg.ReportCaller = false

	logger, err := log.NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	defer logger.Close()

	// 批量添加字段
	ctx := log.WithFields(context.Background(), map[string]any{
		"field1": "value1",
		"field2": "value2",
		"field3": "value3",
		"field4": 123,
		"field5": true,
	})

	logger.InfoContextf(ctx, "Message with multiple fields")

	output := buf.String()

	// 验证所有字段都存在
	expectedFields := []string{"field1", "field2", "field3", "field4", "field5"}
	for _, field := range expectedFields {
		if !strings.Contains(output, `"`+field+`"`) {
			t.Errorf("Field '%s' not found in output", field)
		}
	}
}

// TestLoggerCloseWithFile 测试 Logger 关闭文件资源
// 注意：不使用全局 Logger，因为全局 Logger 是单例，在其他测试中已经初始化
func TestLoggerCloseWithFile(t *testing.T) {
	tmpFile := "/tmp/test_logger_close.log"
	defer func() {
		// 确保清理临时文件
		if _, err := os.Stat(tmpFile); err == nil {
			os.Remove(tmpFile)
		}
	}()

	cfg := log.DefaultConfig()
	cfg.FilePath = tmpFile

	// 创建独立的 Logger 实例
	logger, err := log.NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	logger.Infof("Test message")

	// 关闭 Logger
	err = logger.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// 验证文件存在且可以被删除（说明句柄已正确关闭）
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("Log file should exist after close")
		return
	}

	err = os.Remove(tmpFile)
	if err != nil {
		t.Errorf("Failed to remove log file: %v", err)
	}
}

// TestCallerSkipFramesConfig 测试自定义 CallerSkipFrames 配置
func TestCallerSkipFramesConfig(t *testing.T) {
	buf := &bytes.Buffer{}
	cfg := log.DefaultConfig()
	cfg.Format = log.FormatJSON
	cfg.Output = buf
	cfg.ReportCaller = true
	cfg.CallerSkipFrames = 0 // 使用默认值

	logger, err := log.NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	defer logger.Close()

	logger.Infof("Test caller info")

	output := buf.String()
	t.Logf("Output: %s", output)

	// 验证 file 字段存在且包含 log_test.go
	if !strings.Contains(output, `"file"`) {
		t.Error("file field not found in output")
	}
	if !strings.Contains(output, "log_test.go") {
		t.Error("file field should contain 'log_test.go'")
	}
}
