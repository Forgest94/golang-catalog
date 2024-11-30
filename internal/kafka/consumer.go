package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

const (
	commitInterval = 5000
	timeoutMessage = -1
)

type Handler interface {
	HandlerMessage(message []byte) error
}

type Consumer struct {
	consumer *kafka.Consumer
	handler  Handler
	stop     bool
}

func NewConsumer(handler Handler, address []string, groupId, topic string) (*Consumer, error) {
	cfg := kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(address, ","),
		"group.id":                 groupId,
		"enable.auto.offset.store": false,
		"auto.commit.interval.ms":  commitInterval,
		"security.protocol":        "sasl_plaintext",
		"sasl.mechanisms":          "PLAIN",
		"sasl.username":            os.Getenv("KAFKA_USER"),
		"sasl.password":            os.Getenv("KAFKA_PASSWORD"),
	}

	c, err := kafka.NewConsumer(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating consumer for group %s: %v", groupId, err)
	}

	if err := c.Subscribe(topic, nil); err != nil {
		return nil, fmt.Errorf("error subscribing to topic %s: %v", topic, err)
	}

	return &Consumer{
		consumer: c,
		handler:  handler,
	}, nil
}

func (c *Consumer) Start() {
	for {
		if c.stop {
			break
		}
		kafkaMessage, err := c.consumer.ReadMessage(timeoutMessage)
		if err != nil {
			logrus.Error(err)
		}

		if kafkaMessage.Value == nil {
			continue
		}

		if err := c.handler.HandlerMessage(kafkaMessage.Value); err != nil {
			logrus.Error(err)
			continue
		}

		if _, err := c.consumer.StoreMessage(kafkaMessage); err != nil {
			logrus.Error(err)
			continue
		}
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	if _, err := c.consumer.Commit(); err != nil {
		return err
	}
	return c.consumer.Close()
}
