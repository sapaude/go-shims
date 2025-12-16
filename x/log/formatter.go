package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
)

// OrderedJSONFormatter 实现可控制字段顺序的 JSON 格式化器
// 字段输出顺序：time → level → file → func → 业务字段（按字母序） → msg
type OrderedJSONFormatter struct {
	// TimestampFormat 时间戳格式，默认为 "2006/01/02 15:04:05.000"
	TimestampFormat string

	// PrettyPrint 是否美化输出（带缩进的 JSON）
	PrettyPrint bool

	// FieldOrder 前置字段的顺序，默认为 ["time", "level", "file", "func"]
	// 这些字段会按照指定顺序首先输出
	FieldOrder []string
}

// Format 实现 logrus.Formatter 接口
func (f *OrderedJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 使用 map 临时存储所有字段
	data := make(map[string]interface{})

	// 1. 添加时间戳
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}
	data["time"] = entry.Time.Format(timestampFormat)

	// 2. 添加日志级别
	data["level"] = entry.Level.String()

	// 3. 添加消息
	data["msg"] = entry.Message

	// 4. 添加 entry.Data 中的所有字段（包括 file、func 和业务字段）
	for k, v := range entry.Data {
		data[k] = v
	}

	// 构建有序的 JSON
	buf := &bytes.Buffer{}

	if f.PrettyPrint {
		buf.WriteString("{\n")
	} else {
		buf.WriteString("{")
	}

	// 获取字段顺序配置
	fieldOrder := f.FieldOrder
	if len(fieldOrder) == 0 {
		fieldOrder = []string{"time", "level", "file", "func"}
	}

	// 分类字段
	orderedFields := make([]string, 0) // 按配置顺序的字段
	businessFields := make([]string, 0) // 业务字段（需要排序）
	hasMsg := false

	// 创建已输出字段的 map 用于去重
	outputFields := make(map[string]bool)

	// 首先按照 FieldOrder 输出前置字段
	for _, field := range fieldOrder {
		if val, ok := data[field]; ok {
			orderedFields = append(orderedFields, field)
			outputFields[field] = true
			_ = val // 避免未使用警告
		}
	}

	// 收集剩余的业务字段（排除 msg）
	for field := range data {
		if field == "msg" {
			hasMsg = true
			continue
		}
		if !outputFields[field] {
			businessFields = append(businessFields, field)
		}
	}

	// 业务字段按字母顺序排序
	sort.Strings(businessFields)

	// 合并所有字段：前置字段 + 业务字段 + msg
	allFields := make([]string, 0, len(orderedFields)+len(businessFields)+1)
	allFields = append(allFields, orderedFields...)
	allFields = append(allFields, businessFields...)
	if hasMsg {
		allFields = append(allFields, "msg")
	}

	// 按顺序输出 JSON 字段
	for i, field := range allFields {
		val := data[field]

		if f.PrettyPrint {
			buf.WriteString("  ")
		}

		// 写入字段名
		buf.WriteString(`"`)
		buf.WriteString(field)
		buf.WriteString(`":`)

		if f.PrettyPrint {
			buf.WriteString(" ")
		}

		// 写入字段值（使用 JSON 编码）
		valueBytes, err := json.Marshal(val)
		if err != nil {
			// 如果序列化失败，使用字符串形式
			valueBytes = []byte(fmt.Sprintf(`"%v"`, val))
		}
		buf.Write(valueBytes)

		// 添加逗号（最后一个字段不需要）
		if i < len(allFields)-1 {
			buf.WriteString(",")
		}

		if f.PrettyPrint {
			buf.WriteString("\n")
		}
	}

	if f.PrettyPrint {
		buf.WriteString("}\n")
	} else {
		buf.WriteString("}\n")
	}

	return buf.Bytes(), nil
}
