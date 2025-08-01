package kafkax

import (
    "encoding/json"

    "github.com/pkg/errors"
)

// ConsumerGroupConfig kafka配置
type ConsumerGroupConfig struct {
    BrokersAddress    []string `json:"brokers_address"`
    ConsumerGroupID   string   `json:"consumer_group_id"`   // 消费者组ID
    GroupIDConsistent bool     `json:"group_id_consistent"` // 消费者组ID是否一致，一些场景需要全量消费消息
    Topics            []string `json:"topics"`
    Consumer          struct {
        OffsetInitial int64 `json:"offset_initial"` // OffsetNewest int64 = -1 , OffsetOldest int64 = -2
    } `json:"consumer"`
}

// ParseConsumerGroupConfig 解析kafka消费者组配置
func ParseConsumerGroupConfig(data []byte) (*ConsumerGroupConfig, error) {
    var c ConsumerGroupConfig
    err := json.Unmarshal(data, &c)
    if err != nil {
        return nil, errors.Wrap(err, "failed to parse kafka config")
    }
    return &c, nil
}
