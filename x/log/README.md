# go-shims/x/log

基于 [logrus](https://github.com/sirupsen/logrus) 的 Go 日志库，提供结构化日志、Context 集成、字段排序等功能。

## 特性

- ✅ **结构化日志**：支持 JSON 和文本格式
- ✅ **字段排序**：JSON 输出按照固定顺序（`time → level → file → func → 业务字段 → msg`）
- ✅ **Context 集成**：自动提取 request_id、trace_id、user_id 等字段
- ✅ **调用者信息**：自动记录文件名、行号、函数名
- ✅ **资源管理**：正确关闭文件句柄，避免资源泄漏
- ✅ **动态配置**：运行时修改日志级别、格式、输出
- ✅ **批量操作**：高性能批量添加日志字段
- ✅ **全局单例**：简化全局日志使用

---

## 快速开始

### 安装

```bash
go get github.com/sapaude/go-shims/x/log
```

### 基本用法

```go
package main

import (
    "github.com/sapaude/go-shims/x/log"
    "github.com/sirupsen/logrus"
)

func main() {
    // 方式 1：使用默认全局 Logger
    log.Infof("This is an info message")
    log.Errorf("This is an error: %v", err)

    // 方式 2：显式初始化全局 Logger
    cfg := log.DefaultConfig()
    cfg.Level = logrus.DebugLevel
    cfg.Format = log.FormatJSON

    if err := log.InitGlobalLogger(cfg); err != nil {
        log.Errorf("Failed to init logger: %v", err)
    }

    log.Debugf("Debug message")
    log.Infof("Info message")
}
```

---

## 配置选项

### Config 结构

```go
type Config struct {
    Level            logrus.Level // 日志级别 (Debug, Info, Warn, Error, Fatal)
    Format           LogFormat    // 输出格式 (FormatJSON 或 FormatText)
    Output           io.Writer    // 输出目标 (os.Stdout, os.Stderr, 或自定义 Writer)
    FilePath         string       // 文件路径（设置后输出到文件）
    JSONPretty       bool         // JSON 美化输出
    ReportCaller     bool         // 是否报告调用者信息 (文件、行号、函数名)
    TimestampFormat  string       // 时间戳格式，默认 "2006/01/02 15:04:05.000"
    CallerSkipFrames int          // 调用栈跳过帧数，0 表示使用默认值
    FieldOrder       []string     // JSON 字段顺序（前置字段）
}
```

### 默认配置

```go
cfg := log.DefaultConfig()
// 等价于：
cfg := log.Config{
    Level:            logrus.InfoLevel,
    Format:           log.FormatJSON,
    Output:           os.Stdout,
    ReportCaller:     true,
    TimestampFormat:  "2006/01/02 15:04:05.000",
    CallerSkipFrames: 0,
    FieldOrder:       []string{"time", "level", "file", "func"},
}
```

---

## 字段排序

JSON 格式的日志输出按照以下顺序：

1. **固定顺序字段**（由 `FieldOrder` 配置）：
   - `time` - 时间戳
   - `level` - 日志级别
   - `file` - 文件位置
   - `func` - 函数名

2. **业务字段**（按字母序）：
   - `request_id`, `trace_id`, `user_id`, `span_id`
   - 自定义字段

3. **消息字段**（始终在最后）：
   - `msg` - 日志消息

### 示例输出

```json
{
  "time": "2025/12/16 10:30:45.123",
  "level": "info",
  "file": "/path/to/main.go:42",
  "func": "main()",
  "request_id": "req-12345",
  "trace_id": "trace-xyz",
  "user_id": "user-456",
  "custom_key": "custom_value",
  "msg": "Request processed successfully"
}
```

### 自定义字段顺序

```go
cfg := log.DefaultConfig()
cfg.FieldOrder = []string{"time", "level", "msg"} // 自定义前置字段顺序

logger, _ := log.NewLogger(cfg)
logger.Infof("Custom order message")
```

---

## Context 集成

### 添加预定义字段

```go
import (
    "context"
    "github.com/sapaude/go-shims/x/log"
)

func handleRequest(ctx context.Context) {
    // 添加单个字段
    ctx = log.WithRequestID(ctx, "req-12345")
    ctx = log.WithUserID(ctx, "user-456")
    ctx = log.WithTraceID(ctx, "trace-xyz")
    ctx = log.WithSpanID(ctx, "span-abc")

    // 使用带 Context 的日志方法
    log.InfoContextf(ctx, "Processing request")
    // 输出：{"time":"...","level":"info",...,"request_id":"req-12345","user_id":"user-456","trace_id":"trace-xyz","span_id":"span-abc","msg":"Processing request"}
}
```

### 批量添加字段（性能优化）

```go
// 推荐方式：批量添加（只复制一次 map）
ctx = log.WithFields(context.Background(), map[string]any{
    "order_id": "order-789",
    "product_id": "prod-123",
    "quantity": 5,
    "price": 99.99,
})

log.InfoContextf(ctx, "Order created")
```

### 便捷方法：同时添加 TraceID 和 SpanID

```go
// 等价于 WithTraceID + WithSpanID
ctx = log.WithTracing(ctx, "trace-xyz", "span-abc")

log.InfoContextf(ctx, "Traced operation")
```

---

## 资源管理

### 文件输出时必须关闭

```go
func main() {
    cfg := log.DefaultConfig()
    cfg.FilePath = "/var/log/app.log"

    logger, err := log.NewLogger(cfg)
    if err != nil {
        log.Fatalf("Failed to create logger: %v", err)
    }
    defer logger.Close() // ✅ 必须关闭以释放文件句柄

    logger.Infof("Application started")
}
```

### 全局 Logger 的关闭

```go
func main() {
    cfg := log.DefaultConfig()
    cfg.FilePath = "/var/log/app.log"
    log.InitGlobalLogger(cfg)

    // 在程序退出前关闭
    defer log.CloseGlobalLogger()

    log.Infof("Application running")
}
```

---

## 动态配置

### 运行时修改日志级别

```go
logger := log.GetGlobalLogger()

// 修改为 Debug 级别
logger.SetLevel(logrus.DebugLevel)

// 修改为 Error 级别
logger.SetLevel(logrus.ErrorLevel)
```

### 运行时修改输出格式

```go
logger := log.GetGlobalLogger()

// 切换到文本格式
logger.SetFormatter(log.FormatText)

// 切换到 JSON 格式
logger.SetFormatter(log.FormatJSON)
```

### 运行时修改输出目标

```go
import "os"

logger := log.GetGlobalLogger()

// 输出到文件
file, _ := os.OpenFile("/var/log/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
logger.SetOutput(file)

// 输出到 stderr
logger.SetOutput(os.Stderr)
```

---

## 创建独立的 Logger 实例

```go
// 创建多个 Logger 用于不同模块
func createModuleLogger(module string, level logrus.Level) log.Logger {
    cfg := log.DefaultConfig()
    cfg.Level = level
    cfg.FilePath = fmt.Sprintf("/var/log/%s.log", module)

    logger, err := log.NewLogger(cfg)
    if err != nil {
        log.Fatalf("Failed to create %s logger: %v", module, err)
    }

    return logger
}

func main() {
    dbLogger := createModuleLogger("database", logrus.DebugLevel)
    defer dbLogger.Close()

    apiLogger := createModuleLogger("api", logrus.InfoLevel)
    defer apiLogger.Close()

    dbLogger.Debugf("Database query executed")
    apiLogger.Infof("API request received")
}
```

---

## 与 Kafka 集成示例

参考 `infra/kafkax` 模块的使用：

```go
import (
    "context"
    "github.com/sapaude/go-shims/x/log"
)

func handleKafkaMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
    // 添加消息元数据到 Context
    ctx = log.WithFields(ctx, map[string]any{
        "topic":     msg.Topic,
        "partition": msg.Partition,
        "offset":    msg.Offset,
    })

    log.InfoContextf(ctx, "Processing Kafka message")

    // 业务处理...
    if err := processMessage(msg); err != nil {
        log.ErrorContextf(ctx, "Failed to process message: %v", err)
        return err
    }

    log.InfoContextf(ctx, "Message processed successfully")
    return nil
}
```

---

## 高级特性

### 调整 CallerSkipFrames

如果调用者信息不准确，可以手动调整栈帧跳过数：

```go
cfg := log.DefaultConfig()
cfg.CallerSkipFrames = 10 // 根据实际调用链深度调整

logger, _ := log.NewLogger(cfg)
logger.Infof("Test message")
```

### 自定义时间戳格式

```go
cfg := log.DefaultConfig()
cfg.TimestampFormat = time.RFC3339 // ISO 8601 格式
// 或
cfg.TimestampFormat = "2006-01-02 15:04:05" // 自定义格式

logger, _ := log.NewLogger(cfg)
```

### JSON 美化输出（用于开发调试）

```go
cfg := log.DefaultConfig()
cfg.JSONPretty = true

logger, _ := log.NewLogger(cfg)
logger.Infof("Pretty JSON output")
// 输出：
// {
//   "time": "2025/12/16 10:30:45.123",
//   "level": "info",
//   "file": "/path/to/file.go:42",
//   "func": "main()",
//   "msg": "Pretty JSON output"
// }
```

---

## API 参考

### 全局日志方法

```go
// 标准日志方法
log.Debugf(format string, args ...any)
log.Infof(format string, args ...any)
log.Warnf(format string, args ...any)
log.Errorf(format string, args ...any)
log.Fatalf(format string, args ...any) // 打印后调用 os.Exit(1)

// 带 Context 的日志方法
log.DebugContextf(ctx context.Context, format string, args ...any)
log.InfoContextf(ctx context.Context, format string, args ...any)
log.WarnContextf(ctx context.Context, format string, args ...any)
log.ErrorContextf(ctx context.Context, format string, args ...any)
log.FatalContextf(ctx context.Context, format string, args ...any)
```

### Context 工具方法

```go
// 单字段操作
log.WithRequestID(ctx context.Context, reqID string) context.Context
log.WithUserID(ctx context.Context, userID string) context.Context
log.WithTraceID(ctx context.Context, traceID string) context.Context
log.WithSpanID(ctx context.Context, spanID string) context.Context
log.WithCustomField(ctx context.Context, key string, value any) context.Context

// 批量操作（推荐）
log.WithFields(ctx context.Context, fields map[string]any) context.Context
log.WithTracing(ctx context.Context, traceID, spanID string) context.Context

// 字段提取
log.GetRequestID(ctx context.Context) (string, bool)
log.GetUserID(ctx context.Context) (string, bool)
log.GetTraceID(ctx context.Context) (string, bool)
log.GetSpanID(ctx context.Context) (string, bool)
log.GetCustomFields(ctx context.Context) (log.MetaData, bool)
```

### Logger 接口

```go
type Logger interface {
    // 标准方法
    Debugf(format string, args ...any)
    Infof(format string, args ...any)
    Warnf(format string, args ...any)
    Errorf(format string, args ...any)
    Fatalf(format string, args ...any)

    // 带 Context 方法
    DebugContextf(ctx context.Context, format string, args ...any)
    InfoContextf(ctx context.Context, format string, args ...any)
    WarnContextf(ctx context.Context, format string, args ...any)
    ErrorContextf(ctx context.Context, format string, args ...any)
    FatalContextf(ctx context.Context, format string, args ...any)

    // 动态配置
    SetLevel(level logrus.Level)
    SetOutput(output io.Writer)
    SetFormatter(format LogFormat)

    // 资源管理
    Close() error
}
```

---

## 最佳实践

### 1. 优先使用批量操作

```go
// ❌ 不推荐：多次调用会导致多次 map 复制
ctx = log.WithCustomField(ctx, "key1", "value1")
ctx = log.WithCustomField(ctx, "key2", "value2")
ctx = log.WithCustomField(ctx, "key3", "value3")

// ✅ 推荐：批量添加，只复制一次 map
ctx = log.WithFields(ctx, map[string]any{
    "key1": "value1",
    "key2": "value2",
    "key3": "value3",
})
```

### 2. 文件输出必须关闭

```go
// ❌ 不推荐：忘记关闭会导致文件句柄泄漏
logger, _ := log.NewLogger(cfg)
logger.Infof("message")

// ✅ 推荐：使用 defer 确保关闭
logger, _ := log.NewLogger(cfg)
defer logger.Close()
logger.Infof("message")
```

### 3. 在 HTTP Handler 中使用 Context

```go
func handleHTTP(w http.ResponseWriter, r *http.Request) {
    // 在请求开始时创建带有请求信息的 Context
    ctx := log.WithRequestID(r.Context(), generateRequestID())
    ctx = log.WithFields(ctx, map[string]any{
        "method": r.Method,
        "path":   r.URL.Path,
        "remote": r.RemoteAddr,
    })

    // 传递 Context 到后续处理
    if err := processRequest(ctx, r); err != nil {
        log.ErrorContextf(ctx, "Request failed: %v", err)
        http.Error(w, "Internal Server Error", 500)
        return
    }

    log.InfoContextf(ctx, "Request succeeded")
}
```

### 4. 生产环境推荐配置

```go
cfg := log.DefaultConfig()
cfg.Level = logrus.InfoLevel    // 生产环境使用 Info 级别
cfg.Format = log.FormatJSON      // JSON 格式便于日志收集
cfg.FilePath = "/var/log/app.log"
cfg.ReportCaller = true          // 记录调用者信息便于调试
cfg.JSONPretty = false           // 生产环境不美化（节省空间）

if err := log.InitGlobalLogger(cfg); err != nil {
    log.Fatalf("Failed to init logger: %v", err)
}
defer log.CloseGlobalLogger()
```

---

## 性能优化建议

1. **使用 WithFields 批量添加字段**：减少 Context 分配和 map 复制
2. **生产环境关闭 JSONPretty**：减少输出体积
3. **合理设置日志级别**：避免输出过多 Debug 日志
4. **使用异步日志收集**：避免日志 I/O 阻塞主流程

---

## 依赖

- [logrus](https://github.com/sirupsen/logrus) v1.9.3

---

## License

MIT

---

## 版本历史

### v0.2.9（当前版本）
- ✅ 添加自定义字段排序支持
- ✅ 修复文件句柄资源泄漏问题
- ✅ 添加 Close() 方法
- ✅ 改进错误处理（InitGlobalLogger 返回 error）
- ✅ 添加批量字段操作（WithFields）
- ✅ 添加便捷方法（WithTracing）
- ✅ 优化 CallerHook 配置
- ✅ 完善测试套件

### v0.2.8
- 初始版本
