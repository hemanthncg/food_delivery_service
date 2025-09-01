package common

import (
    "context"
    "github.com/segmentio/kafka-go"
    "log"
    "os"
    "strings"
)

var (
    kafkaBrokers = initBrokers()
)

// GetKafkaBrokers returns configured brokers
func GetKafkaBrokers() []string {
    return kafkaBrokers
}

func initBrokers() []string {
    env := os.Getenv("KAFKA_BROKERS")
    if env == "" {
        return []string{"kafka:9092"}
    }
    parts := strings.Split(env, ",")
    var brokers []string
    for _, p := range parts {
        v := strings.TrimSpace(p)
        if v != "" {
            brokers = append(brokers, v)
        }
    }
    if len(brokers) == 0 {
        brokers = []string{"kafka:9092"}
    }
    return brokers
}

func NewKafkaWriter(topic string) *kafka.Writer {
    return &kafka.Writer{
        Addr:     kafka.TCP(kafkaBrokers...),
        Topic:    topic,
        Balancer: &kafka.LeastBytes{},
    }
}

func PublishKafka(topic string, key, value []byte) error {
    w := NewKafkaWriter(topic)
    defer w.Close()
    return w.WriteMessages(context.Background(),
        kafka.Message{Key: key, Value: value},
    )
}

func NewKafkaReader(topic, groupID string) *kafka.Reader {
    return kafka.NewReader(kafka.ReaderConfig{
        Brokers:  kafkaBrokers,
        GroupID:  groupID,
        Topic:    topic,
        MinBytes: 1,
        MaxBytes: 10e6,
    })
}

func ConsumeKafka(topic, groupID string, handler func([]byte)) {
    r := NewKafkaReader(topic, groupID)
    defer r.Close()
    for {
        m, err := r.ReadMessage(context.Background())
        if err != nil {
            log.Println("kafka read error:", err)
            continue
        }
        handler(m.Value)
    }
}
