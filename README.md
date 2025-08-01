# go-shims

go-shims 是一个用于Go项目开发的工具库，提供字符串、数值、文件、路径等常见操作的便捷函数，帮助开发者简化代码，提高开发效率。

## 安装使用

### shim - 垫片函数

```shell
# 业务垫片库
go get -u github.com/sapaude/go-shims/shim
```

### 垫片支持功能模块

- 字符串处理
- 数值转换
- 文件操作
- 路径处理

## X库使用

```shell
# 业务基础扩展库，例如日志库
go get -u github.com/sapaude/go-shims/x/log
```

### 支持扩展能力

- `x`目录下后续会支持各类小型包，例如日志、数据库、AI组件等

## Infrastructure基础设施库（例如Kafka\DB\COS等）

- `kafkax`: 支持快速kafka包使用

