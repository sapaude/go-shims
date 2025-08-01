# github.com/sapaude/go-shims/infra/kafkax

## 极简使用

```shell
go get -u github.com/sapaude/go-shims/infra/kafkax
```

1. 解析kafka配置
2. 初始化消费者组app应用
3. 初始化消费者处理方法handler，实现`sarama.ConsumerGroupHandler`接口

## 示例

```go
type exampleConsumerApp struct {
    consumerApp *kafkax.ConsumerGroupApp
}

func newExampleConsumerApp() *exampleConsumerApp {
    data, err := os.ReadFile("./consumer_config.json")
    if err != nil {
        return nil
    }
    cfg, err := kafkax.ParseConsumerGroupConfig(data)
    if err != nil {
        fmt.Errorf("parse consumer_config.json fail: %s", err.Error())
    }

    consumerApp := kafkax.NewConsumerApp(cfg, newExampleConsumer())
    return &exampleConsumerApp{consumerApp: consumerApp}
}

// StartConsume 启动消费
func (e *exampleConsumerApp) StartConsume() {
    e.consumerApp.StartConsume()
}
```