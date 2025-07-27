package shim

import (
    "regexp"
    "strings"
)

// 定义正则表达式，使用 MustCompile 确保在程序启动时编译成功
var (
    // reMarkdownJsonBlock 匹配 ```json\n...``` 形式的JSON代码块
    // (?s) 使得 . 匹配包括换行符在内的所有字符
    // (.*?) 非贪婪匹配，确保只匹配到第一个 ``` 结束符
    reMarkdownJsonBlock = regexp.MustCompile("(?s)```json\n(.*?)```")

    // reGenericMarkdownBlock 匹配 ```\n...``` 形式的通用代码块
    reGenericMarkdownBlock = regexp.MustCompile("(?s)```\n(.*?)```")

    // reInlineMarkdownJsonBlock 匹配 ```json{...}``` 这种不规范的，json紧跟在```后面的情况
    reInlineMarkdownJsonBlock = regexp.MustCompile("(?s)```json({.*})```")

    // reInlineGenericMarkdownBlock 匹配 ```{...}``` 这种不规范的，没有json标记的情况
    reInlineGenericMarkdownBlock = regexp.MustCompile("(?s)```({.*})```")

    // reLeadingTrailingBraces 匹配以 { 或 [ 开始，以 } 或 ] 结束的字符串
    // 这是一个更通用的匹配，用于捕获纯JSON或被文本包裹的JSON
    reLeadingTrailingBraces = regexp.MustCompile(`^[{\[].*?[}\]]$`)
)

// ExtractPotentialJSON 尝试从原始响应字符串中提取潜在的JSON内容。
// 它采用多策略匹配，以提高鲁棒性。
func ExtractPotentialJSON(rawResponse string) string {
    // 1. 尝试匹配 ```json\n...``` 形式的JSON代码块
    if matches := reMarkdownJsonBlock.FindStringSubmatch(rawResponse); len(matches) > 1 {
        return cleanJSONString(matches[1])
    }

    // 2. 尝试匹配 ```\n...``` 形式的通用代码块
    if matches := reGenericMarkdownBlock.FindStringSubmatch(rawResponse); len(matches) > 1 {
        // 进一步验证提取的内容是否看起来像JSON (以 { 或 [ 开头)
        trimmed := strings.TrimSpace(matches[1])
        if (strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) ||
            (strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]")) {
            return cleanJSONString(trimmed)
        }
    }

    // 3. 尝试匹配 ```json{...}``` 这种不规范的内联JSON代码块
    if matches := reInlineMarkdownJsonBlock.FindStringSubmatch(rawResponse); len(matches) > 1 {
        return cleanJSONString(matches[1])
    }

    // 4. 尝试匹配 ```{...}``` 这种不规范的内联通用代码块
    if matches := reInlineGenericMarkdownBlock.FindStringSubmatch(rawResponse); len(matches) > 1 {
        return cleanJSONString(matches[1])
    }

    // 5. 如果没有找到代码块，尝试查找第一个 { 或 [ 到最后一个 } 或 ]
    // 这种方法更通用，但可能包含JSON前后的额外文本，需要后续的json.Unmarshal来验证
    if matches := reLeadingTrailingBraces.FindStringSubmatch(rawResponse); len(matches) > 0 {
        // FindStringSubmatch 返回的是整个匹配，这里我们直接用它
        return cleanJSONString(matches[0])
    }

    // 6. 最后，尝试去除常见的LLM前缀/后缀和Markdown标记，并清理空白
    // 这种方法最不精确，但作为最后的尝试
    cleaned := strings.ReplaceAll(rawResponse, "```json", "")
    cleaned = strings.ReplaceAll(cleaned, "```", "")
    cleaned = strings.ReplaceAll(cleaned, "Here is the JSON:", "")
    cleaned = strings.ReplaceAll(cleaned, "```json\n", "") // 再次清理，以防前面没匹配到完整块
    cleaned = strings.ReplaceAll(cleaned, "```\n", "")

    // 尝试找到第一个 { 或 [ 和最后一个 } 或 ]
    firstBrace := strings.IndexAny(cleaned, "{[")
    lastBrace := strings.LastIndexAny(cleaned, "}]")

    if firstBrace != -1 && lastBrace != -1 && lastBrace > firstBrace {
        potentialJSON := cleaned[firstBrace : lastBrace+1]
        return cleanJSONString(potentialJSON)
    }

    return "" // 如果所有尝试都失败，返回空字符串
}

// cleanJSONString 对提取到的潜在JSON字符串进行一些清理
func cleanJSONString(s string) string {
    s = strings.TrimSpace(s)
    // 移除JSON对象或数组末尾的逗号（如果存在且不合法）
    // 这是一个常见的LLM生成错误，但需要小心处理，避免误删合法逗号
    // 简单的策略是移除紧跟在 } 或 ] 之前的逗号
    if len(s) > 1 {
        lastChar := s[len(s)-1]
        secondLastChar := s[len(s)-2]
        if (lastChar == '}' || lastChar == ']') && secondLastChar == ',' {
            s = s[:len(s)-2] + string(lastChar) // 移除逗号，保留 } 或 ]
        }
    }
    return s
}
