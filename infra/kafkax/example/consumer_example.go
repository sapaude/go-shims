package example

import (
    "fmt"
    "os"

    "github.com/sapaude/go-shims/infra/kafkax"
)

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
