package kafkax

import (
    "context"
    "fmt"
    "time"

    "github.com/IBM/sarama"
    "github.com/sapaude/go-shims/shim"
    "github.com/sapaude/go-shims/x/log"
)

// ConsumerGroupApp 注入消息处理函数，同时创建一个消费者组处理应用
type ConsumerGroupApp struct {
    config               *ConsumerGroupConfig
    consumerGroupHandler sarama.ConsumerGroupHandler
}

// NewConsumerApp 创建一个通用的ConsumerApp
func NewConsumerApp(config *ConsumerGroupConfig, handler sarama.ConsumerGroupHandler) *ConsumerGroupApp {
    return &ConsumerGroupApp{
        config:               config,
        consumerGroupHandler: handler,
    }
}

// StartConsume 启动消费
func (app *ConsumerGroupApp) StartConsume() {
    // 消费配置
    kfkCfg := app.config
    groupID := app.GetConsumerGroupID() // 消费者组名
    topics := kfkCfg.Topics             // 主题

    // 补充消费者偏移量配置
    saramaConfig := sarama.NewConfig()
    saramaConfig.Consumer.Offsets.Initial = kfkCfg.Consumer.OffsetInitial

    // 创建消费者组实例
    log.Infof(
        "[MQConsume] New Consume GroupId(%s), Kafka Config: %s",
        groupID,
        shim.ToJsonString(kfkCfg, false),
    )
    cgroup, err := sarama.NewConsumerGroup(kfkCfg.BrokersAddress, groupID, saramaConfig)
    if err != nil {
        log.Fatalf("create Kafka consumer group got err: %s", err)
    }

    // 加入消费者集群，消费指定的topics，并启动消费
    for {
        ctx := context.Background()
        if err = cgroup.Consume(ctx, topics, app.consumerGroupHandler); err != nil {
            log.ErrorContextf(ctx, "[MQConsume] Error on Consume with err: %s", err)
        }
        time.Sleep(3 * time.Second)
    }
}

// GetConsumerGroupID 消费者组名
func (app *ConsumerGroupApp) GetConsumerGroupID() string {
    // 分组ID一致
    if app.config.GroupIDConsistent {
        return app.config.ConsumerGroupID
    }

    // 分组ID不一致
    // return fmt.Sprintf("%s-%s", app.config.ConsumerGroupID, shim.GetLocalIP())
    return fmt.Sprintf("%s-%s", app.config.ConsumerGroupID, shim.GenRandomLengthStr(8))
}
