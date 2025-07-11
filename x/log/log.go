package log

import (
    "context"
    "fmt"
    "log/slog"
    "os"
)

const (
    LevelTrace = slog.LevelDebug - 4 // 比 Debug 更低
    LevelFatal = slog.LevelError + 4
)

func Default() *slog.Logger {
    return slog.Default()
}

// init 函数用于在包加载时设置默认的 slog Logger。
// 这是一个常见的模式，确保在任何日志调用发生之前，Logger 已经被配置。
func init() {
    // 使用自定义的 PrettySourceJSONHandler
    slog.SetDefault(slog.New(NewPrettySourceJSONHandler(os.Stdout, &slog.HandlerOptions{
        AddSource: false, // 仍然需要设置为 true，以便 slog 填充 Record.PC
        Level:     slog.LevelDebug,
    })))
}

// SetDefaultLogger 允许外部设置默认的 slog.Logger 实例。
// 这在测试或需要动态切换日志配置时非常有用。
func SetDefaultLogger(l *slog.Logger) {
    slog.SetDefault(l)
}

// DebugContextf 记录 DEBUG 级别的格式化日志，并包含上下文。
func DebugContextf(ctx context.Context, format string, args ...any) {
    logger := slog.Default()
    if !logger.Enabled(ctx, slog.LevelDebug) {
        return
    }
    logger.DebugContext(ctx, fmt.Sprintf(format, args...))
}

// InfoContextf 记录 INFO 级别的格式化日志，并包含上下文。
func InfoContextf(ctx context.Context, format string, args ...any) {
    // 使用 slog.Default() 获取当前默认的 Logger
    logger := slog.Default()
    if !logger.Enabled(ctx, slog.LevelInfo) {
        return
    }
    logger.InfoContext(ctx, fmt.Sprintf(format, args...))
}

// WarnContextf 记录 WARN 级别的格式化日志，并包含上下文。
func WarnContextf(ctx context.Context, format string, args ...any) {
    logger := slog.Default()
    if !logger.Enabled(ctx, slog.LevelWarn) {
        return
    }
    logger.WarnContext(ctx, fmt.Sprintf(format, args...))
}

// ErrorContextf 记录 ERROR 级别的格式化日志，并包含上下文。
func ErrorContextf(ctx context.Context, format string, args ...any) {
    logger := slog.Default()
    if !logger.Enabled(ctx, slog.LevelError) {
        return
    }
    logger.ErrorContext(ctx, fmt.Sprintf(format, args...))
}

func FatalContextf(ctx context.Context, format string, args ...any) {
    logger := slog.Default()
    if !logger.Enabled(ctx, LevelFatal) {
        return
    }
    logger.Log(ctx, LevelFatal, fmt.Sprintf(format, args...))
    os.Exit(1)
}

func Debugf(format string, args ...any) {
    logger := slog.Default()
    logger.Debug(fmt.Sprintf(format, args...))
}

func Infof(format string, args ...any) {
    logger := slog.Default()
    logger.Info(fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...any) {
    logger := slog.Default()
    logger.Warn(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...any) {
    logger := slog.Default()
    logger.Error(fmt.Sprintf(format, args...))
}

func Fatalf(format string, args ...any) {
    logger := slog.Default()
    logger.Log(context.Background(), LevelFatal, fmt.Sprintf(format, args...))
    os.Exit(1)
}
