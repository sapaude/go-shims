package log

import (
    "io"
    "os"

    "github.com/sirupsen/logrus"
)

// LogFormat 定义日志输出格式
type LogFormat string

const (
    FormatText LogFormat = "text"
    FormatJSON LogFormat = "json"
)

// Config 定义日志库的配置参数
type Config struct {
    Level            logrus.Level // 日志级别
    Format           LogFormat    // 日志输出格式 (text/json)
    Output           io.Writer    // 日志输出目标 (例如 os.Stdout, 文件)
    FilePath         string       // 如果输出到文件，指定文件路径
    JSONPretty       bool         // JSON美化输出
    ReportCaller     bool         // 是否报告调用者信息 (文件, 行号, 函数名)
    TimestampFormat  string       // 时间戳格式，默认为 "2006/01/02 15:04:05.000"
    CallerSkipFrames int          // 调用者栈帧跳过数，0 表示使用默认值
    FieldOrder       []string     // JSON 字段输出顺序，仅对 JSON 格式生效
}

// DefaultConfig 返回一个默认的日志配置
func DefaultConfig() Config {
    return Config{
        Level:            logrus.InfoLevel,
        Format:           FormatJSON,
        Output:           os.Stdout,
        FilePath:         "",
        JSONPretty:       false,
        ReportCaller:     true, // 默认开启调用者信息
        TimestampFormat:  "2006/01/02 15:04:05.000",
        CallerSkipFrames: 0,                  // 0 表示使用默认值
        FieldOrder:       DefaultFieldOrder(), // 使用默认字段顺序
    }
}

// DefaultFieldOrder 返回默认的 JSON 字段输出顺序
// 顺序：time → level → file → func → 业务字段（按字母序） → msg
func DefaultFieldOrder() []string {
    return []string{"time", "level", "file", "func"}
}
