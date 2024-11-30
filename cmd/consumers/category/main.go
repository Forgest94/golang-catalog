package main

import (
	"catalog/internal/handlers/kafka/categoryHandler"
	"catalog/internal/kafka"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	topic = "categories"
	group
)

var address = strings.Split(os.Getenv("KAFKA_HOSTS"), ",")

func main() {
	h := categoryHandler.NewHandler()
	c, err := kafka.NewConsumer(h, address, group, topic)
	if err != nil {
		logrus.Fatal(err)
	}
	go c.Start()

	c2, err := kafka.NewConsumer(h, address, group, topic)
	if err != nil {
		logrus.Fatal(err)
	}
	go c2.Start()

	c3, err := kafka.NewConsumer(h, address, group, topic)
	if err != nil {
		logrus.Fatal(err)
	}
	go c3.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	logrus.Fatal(c.Stop(), c2.Stop(), c3.Stop())
}
