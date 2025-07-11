package log

import (
    "context"
    "log/slog"
    "os"
    "testing"
)

func TestLogContextf(t *testing.T) {

    // 设置一个 JSON Handler 作为默认 Logger
    slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        AddSource: true,            // 添加源文件信息
        Level:     slog.LevelDebug, // 设置最低日志级别为 Debug
    })))

    // 使用包装后的 Logger
    ctx := context.Background()

    InfoContextf(ctx, "用户 %s 登录成功，ID 为 %d", "Alice", 123)
    ErrorContextf(ctx, "处理请求失败，错误码：%d，详情：%s", 500, "数据库连接超时")
    DebugContextf(ctx, "调试信息：变量 x = %v", map[string]int{"a": 1, "b": 2})
    WarnContextf(ctx, "缓存命中率低，当前为 %.2f%%", 85.5)
    FatalContextf(ctx, "Fatal err")

    // 也可以添加额外的属性
    // Default().InfoContext(ctx, "用户操作",
    //     slog.String("user", "Bob"),
    //     slog.Int("action_id", 456),
    //     slog.String("message", fmt.Sprintf("用户 %s 进行了操作", "Bob")),
    // )
}

func TestLogf(t *testing.T) {
    Debugf("type: %s %s", "debug", "1")
    Infof("type: %s %s", "info", "2")
    Warnf("type: %s %s", "warn", "3")
    Errorf("type: %s %s", "errro", "4")
    Fatalf("type: %s %s", "fatal", "5")
    t.Logf("not run")
}
