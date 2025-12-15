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
fmt.Errorf("parse consumer_config.json fail: %s", err.BizError())
}

consumerApp := kafkax.NewConsumerApp(cfg, newExampleConsumer())
return &exampleConsumerApp{consumerApp: consumerApp}
}

// StartConsume 启动消费
func (e *exampleConsumerApp) StartConsume() {
e.consumerApp.StartConsume()
}
```

## 基础设施安装

### Mac 安装 kafka本地实例

1. 使用安装`Colima`容器服务，用于拉取Docker镜像
2. 结合`docker-compose.yaml`项目，快速在Mac本地启动一个Kafka本地实例：https://github.com/sapaude/docker-infras/tree/main/mq/kafka
3. 配置了kafka-ui管理kafka topic、消费者等信息：https://github.com/provectus/kafka-ui

```yaml
services:
  kafka:
    image: bitnami/kafka:3.7
    container_name: kafka
    ports:
      - "9092:9092"
    environment:
      - KAFKA_ENABLE_KRAFT=yes
      - KAFKA_CFG_PROCESS_ROLES=controller,broker
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://10.11.12.33:9092
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@127.0.0.1:9093
      - KAFKA_CFG_NODE_ID=1
    volumes:
      - kafka_data:/bitnami/kafka
    networks:
      - kafka_net
    healthcheck:
      test: ["CMD-SHELL", "kafka-topics.sh --bootstrap-server localhost:9092 --list || exit 1"]
      interval: 10s
      timeout: 30s
      retries: 15
      start_period: 40s

  kafka-ui:
    container_name: kafka-ui
    ports:
      - "127.0.0.1:18081:8080"
    image: provectuslabs/kafka-ui:latest
    environment:
      KAFKA_CLUSTERS_0_NAME: kafka
      KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka:9092
    depends_on:
      kafka:
        condition: service_healthy
    networks:
      - kafka_net

volumes:
  kafka_data:

networks:
  kafka_net:
```