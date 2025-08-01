package example

import (
    "context"
    "fmt"

    "github.com/IBM/sarama"
    "github.com/sapaude/go-shims/x/log"
)

type exampleConsumer struct {
}

func newExampleConsumer() *exampleConsumer {
    return &exampleConsumer{}
}

// Setup 方法在每次分区分配（rebalance）完成后都会被调用
func (c *exampleConsumer) Setup(session sarama.ConsumerGroupSession) error {
    fmt.Printf("[MQConsume] Setup: %v", session)
    return nil
}

// Cleanup 方法在每次分区分配（rebalance）完成后都会被调用
func (c *exampleConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
    fmt.Printf("[MQConsume] Cleanup: %v", session)
    return nil
}

// ConsumeClaim 负责消费单一分区的协程，不断从分区中获取消息进行消费
func (c *exampleConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    fmt.Printf("[MQConsume] ConsumeClaim: %v", session)
    for consumeMsg := range claim.Messages() {
        err := c.MsgHandler(context.Background(), consumeMsg)
        if err != nil {
            fmt.Errorf("fn[srv.consumeMsgHandle] SSEMsgSender sessId[%s] consume msg got err: %s", consumeMsg.Key, err.Error())
            continue
        }
        session.MarkMessage(consumeMsg, "")
    }
    return nil
}

func (c *exampleConsumer) MsgHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
    log.Infof("[MQConsume] MsgHandler: %v", msg)

    return nil
}
