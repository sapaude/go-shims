package log

import (
    "context"
    "fmt"
    "io"
    "log/slog"
    "runtime"
)

// PrettySourceJSONHandler 是一个包装了 slog.JSONHandler 的自定义 Handler。
// 它修改了 source 字段的输出格式。
type PrettySourceJSONHandler struct {
    slog.Handler
}

// NewPrettySourceJSONHandler 创建并返回一个 PrettySourceJSONHandler 实例。
func NewPrettySourceJSONHandler(w io.Writer, opts *slog.HandlerOptions) *PrettySourceJSONHandler {
    // 确保 AddSource 选项为 true，这样 slog 才会填充 Record.PC
    if opts == nil {
        opts = &slog.HandlerOptions{}
    }
    // opts.AddSource = false // 强制开启源信息
    return &PrettySourceJSONHandler{slog.NewJSONHandler(w, opts)}
}

// Handle 方法实现了 slog.Handler 接口。
func (h *PrettySourceJSONHandler) Handle(ctx context.Context, r slog.Record) error {
    // 如果 Record 中有 Source 字段，我们修改它
    if r.PC != 0 {
        fs := runtime.CallersFrames([]uintptr{r.PC})
        f, _ := fs.Next()

        // 提取文件名和行号，并格式化
        // shortFile := filepath.Base(f.File)
        sourceStr := fmt.Sprintf("%s:%d", f.File, f.Line)

        // 移除原始的 source 属性，添加我们自定义的
        // 注意：slog.Record 的 Attrs 是一个切片，直接修改可能影响性能或导致意外行为。
        // 更安全的方式是创建一个新的 Record 或在 AddAttrs 中覆盖。
        // 这里我们直接在 Record 上操作，因为我们知道 slog 内部如何处理。
        // 遍历 Record 的属性，找到并替换 source 属性
        var newAttrs []slog.Attr
        foundSource := false
        r.Attrs(func(a slog.Attr) bool {
            if a.Key == slog.SourceKey {
                newAttrs = append(newAttrs, slog.String(slog.SourceKey, sourceStr))
                foundSource = true
            } else {
                newAttrs = append(newAttrs, a)
            }
            return true // 继续遍历
        })

        // 如果没有找到原始的 source 属性（理论上不会发生，因为 AddSource=true），则添加
        if !foundSource {
            newAttrs = append(newAttrs, slog.String(slog.SourceKey, sourceStr))
        }

        // 创建一个新的 Record，包含修改后的属性
        // 注意：slog.Record 是一个结构体，直接修改其内部属性可能不被推荐。
        // 更健壮的方式是创建一个新的 Record，但 slog.Record 没有公共的构造函数来复制所有字段。
        // 这里的做法是利用 slog.Record 的 Attrs 方法来迭代和修改。
        // 实际上，slog.Record 的 Attrs 方法是只读的，不能直接修改。
        // 我们需要通过 WithAttrs 来添加或覆盖属性。
        // 重新构建 Record 的属性列表
        // 这是一个更安全的做法，虽然可能略有性能开销
        r2 := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
        r2.AddAttrs(newAttrs...) // 添加所有属性，包括修改后的 source

        // 将修改后的 Record 传递给底层的 Handler
        return h.Handler.Handle(ctx, r2)
    }

    // 如果没有 Source 字段，则直接传递给底层的 Handler
    return h.Handler.Handle(ctx, r)
}

// WithAttrs 方法实现了 slog.Handler 接口。
func (h *PrettySourceJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    return &PrettySourceJSONHandler{h.Handler.WithAttrs(attrs)}
}

// WithGroup 方法实现了 slog.Handler 接口。
func (h *PrettySourceJSONHandler) WithGroup(name string) slog.Handler {
    return &PrettySourceJSONHandler{h.Handler.WithGroup(name)}
}

// Enabled 方法实现了 slog.Handler 接口。
func (h *PrettySourceJSONHandler) Enabled(ctx context.Context, level slog.Level) bool {
    return h.Handler.Enabled(ctx, level)
}
