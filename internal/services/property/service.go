package property

import (
	"catalog/internal/models/property"
	"encoding/json"
	"fmt"
)

func GetById(id string) (*property.Property, error) {
	p, _ := property.NewModel()
	return p.GetById(id)
}

func LoadMessageConsumerV1(message []byte) error {
	var receivedMessage ConsumerMessageV1
	err := json.Unmarshal(message, &receivedMessage)
	if err != nil {
		return fmt.Errorf("could not unmarshal property message v1: %w", err)
	}

	if receivedMessage.Payload.Id == "" {
		return fmt.Errorf("invalid property message v1")
	}

	p, _ := property.NewModel()
	propertyUpdateFields := property.Property{
		Id:                receivedMessage.Payload.Id,
		IsActive:          receivedMessage.Payload.IsActive,
		Name:              receivedMessage.Payload.Name,
		Code:              receivedMessage.Payload.Code,
		Hint:              receivedMessage.Payload.Hint,
		Type:              receivedMessage.Payload.Type,
		ShowFilter:        receivedMessage.Payload.ShowFilter,
		ShowProductList:   receivedMessage.Payload.ShowProductList,
		ShowDetailProduct: receivedMessage.Payload.ShowDetailProduct,
	}

	getProperty, err := p.GetById(propertyUpdateFields.Id)
	switch receivedMessage.Event {
	case "add":
		if err != nil {
			if err := p.Add(propertyUpdateFields); err != nil {
				return fmt.Errorf("could not add property: %w", err)
			}
		} else {
			return fmt.Errorf("could not add property becouse exist property %d", propertyUpdateFields.Id)
		}
	case "update":
		if err != nil {
			return fmt.Errorf("could not get product: %w", err)
		}
		if err := p.Update(getProperty.Id, propertyUpdateFields); err != nil {
			return fmt.Errorf("could not update property: %w", err)
		}
	case "delete":
		if err != nil {
			return fmt.Errorf("could not get property: %w", err)
		}
		if err := p.Delete(getProperty.Id); err != nil {
			return fmt.Errorf("could not delete property: %w", err)
		}
	default:
		return fmt.Errorf("unknown event: %s", receivedMessage.Event)
	}

	return nil
}
