package kafkax

import (
    "encoding/json"
    "testing"
)

func TestLoggingConfig(t *testing.T) {
    c := &ConsumerGroupConfig{
        BrokersAddress:    []string{"127.0.0.1:9092"},
        ConsumerGroupID:   "consumer_group.sse_push",
        GroupIDConsistent: false,
        Topics:            []string{"test-topic"},
        Consumer: struct {
            OffsetInitial int64 `json:"offset_initial"`
        }{
            OffsetInitial: -1,
        },
    }
    s, err := json.Marshal(c)
    if err != nil {
        return
    }
    t.Log(string(s))
}
