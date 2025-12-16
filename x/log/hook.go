package log

import (
    "fmt"
    "runtime"
    "strings"

    "github.com/sirupsen/logrus"
)

const (
    // CallerFileFieldKey 是日志中存储调用者信息的字段名
    CallerFileFieldKey = "file"
    CallerFuncFieldKey = "func"

    // DefaultCallerSkipFrames 是默认的栈帧跳过数
    // 这个值需要根据调用链深度调整：
    // 调用链：用户代码 → log.Infof → GetGlobalLogger → LogrusLogger.Infof → logrus.Infof → logrus.Entry.logf → hook.Fire → runtime.Caller
    // 需要跳过：hook.Fire(0) → logrus.Entry.logf(1) → logrus.log(2) → logrus.Infof(3) → LogrusLogger.Infof(4) → (可能的包装层) → 到达用户代码
    DefaultCallerSkipFrames = 10
)

// CallerHook 是一个 Logrus Hook，用于添加调用者信息（文件、行号、函数名）
type CallerHook struct {
    // SkipFrames 决定向上跳过多少个栈帧来找到真正的调用者
    // 默认情况下，我们需要跳过 Logrus 内部调用和我们自己的封装层
    SkipFrames int
}

// NewCallerHook 创建一个新的 CallerHook 实例
// skipFrames: 栈帧跳过数，如果 <= 0 则使用 DefaultCallerSkipFrames
func NewCallerHook(skipFrames int) *CallerHook {
    if skipFrames <= 0 {
        skipFrames = DefaultCallerSkipFrames
    }
    return &CallerHook{
        SkipFrames: skipFrames,
    }
}

// Levels 返回 Hook 应该触发的日志级别
func (hook *CallerHook) Levels() []logrus.Level {
    return logrus.AllLevels
}

// Fire 在日志事件发生时被调用
func (hook *CallerHook) Fire(entry *logrus.Entry) error {
    // 向上跳过 hook.Fire, logrus.Entry.log, my_logger.Logger 方法, 以及 Logrus 内部的调用
    // 具体的跳过帧数可能需要根据实际封装层级进行微调
    pc, file, line, ok := runtime.Caller(hook.SkipFrames)
    if !ok {
        return nil
    }

    funcName := runtime.FuncForPC(pc).Name()
    // 简化函数名，去除包路径
    lastSlash := strings.LastIndex(funcName, "/")
    if lastSlash != -1 {
        funcName = funcName[lastSlash+1:]
    }
    lastDot := strings.LastIndex(funcName, ".")
    if lastDot != -1 {
        funcName = funcName[lastDot+1:]
    }

    // 格式化调用者信息
    entry.Data[CallerFileFieldKey] = fmt.Sprintf("%s:%d", file, line)
    entry.Data[CallerFuncFieldKey] = fmt.Sprintf("%s()", funcName)
    return nil

}
