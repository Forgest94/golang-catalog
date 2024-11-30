package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"os"
	"strings"
)

const flushTimeout = 5000

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(address []string, groupId string) (*Producer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
		"group.id":          groupId,
		"security.protocol": "sasl_plaintext",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     os.Getenv("KAFKA_USER"),
		"sasl.password":     os.Getenv("KAFKA_PASSWORD"),
	}

	p, err := kafka.NewProducer(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating kafka producer: %v", err)
	}

	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(topic, key, msg string) error {
	kafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: []byte(msg),
	}

	kafkaChain := make(chan kafka.Event, 1024)
	if err := p.producer.Produce(kafkaMessage, kafkaChain); err != nil {
		return fmt.Errorf("error sending message to kafka: %v", err)
	}

	e := <-kafkaChain
	switch e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return e.(kafka.Error)
	default:
		return fmt.Errorf("unknown message type: %v", e)
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
