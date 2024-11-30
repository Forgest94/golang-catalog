package kafka

type BaseConsumerMessage struct {
	Subject string `json:"subject"`
	Version string `json:"version"`
}
