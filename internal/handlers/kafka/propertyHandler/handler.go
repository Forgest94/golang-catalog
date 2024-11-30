package propertyHandler

import (
	"catalog/internal/handlers/kafka"
	"catalog/internal/services/property"
	"encoding/json"
	"fmt"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandlerMessage(message []byte) error {
	var receivedMessage kafka.BaseConsumerMessage
	err := json.Unmarshal(message, &receivedMessage)
	if err != nil {
		return fmt.Errorf("could not unmarshal property message: %w", err)
	}

	if receivedMessage.Subject != "bd.property.service" {
		return fmt.Errorf("invalid subject")
	}

	switch receivedMessage.Version {
	case "1.0.0":
		if err := property.LoadMessageConsumerV1(message); err != nil {
			return fmt.Errorf("could not load message consumer: %w", err)
		}
	default:
		return fmt.Errorf("invalid version")
	}

	return nil
}
