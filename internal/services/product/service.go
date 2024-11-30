package product

import (
	"catalog/internal/models/product"
	"catalog/internal/services/property"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

func LoadMessageConsumerV1(message []byte) error {
	var receivedMessage ConsumerMessageV1
	err := json.Unmarshal(message, &receivedMessage)
	if err != nil {
		return fmt.Errorf("could not unmarshal product message v1: %w", err)
	}

	if receivedMessage.Payload.Id == "" {
		return fmt.Errorf("invalid product message v1")
	}

	p, _ := product.NewModel()
	productUpdateFields := product.Product{
		Id:                 receivedMessage.Payload.Id,
		ParentId:           receivedMessage.Payload.ParentId,
		IsActive:           receivedMessage.Payload.IsActive,
		Name:               receivedMessage.Payload.Name,
		Code:               receivedMessage.Payload.Code,
		Description:        receivedMessage.Payload.Description,
		PreviewDescription: receivedMessage.Payload.PreviewDescription,
		BasePrice:          receivedMessage.Payload.BasePrice,
		Price:              receivedMessage.Payload.Price,
		Categories:         receivedMessage.Payload.Categories,
	}

	if len(productUpdateFields.Categories) == 0 {
		logrus.Errorf("could not find any categories")
		productUpdateFields.IsActive = false
	}

	if len(receivedMessage.Payload.Properties) > 0 {
		createProperties := make(map[string]struct {
			Id     string `json:"id"`
			Values []struct {
				String  string  `json:"string,omitempty"`
				Integer int     `json:"integer,omitempty"`
				Float   float32 `json:"float,omitempty"`
				Boolean bool    `json:"boolean,omitempty"`
			} `json:"values"`
		}, len(receivedMessage.Payload.Properties))
		for _, payloadProperty := range receivedMessage.Payload.Properties {
			getProperty, err := property.GetById(payloadProperty.Id)
			if err != nil {
				logrus.Errorf("could not get property: %w", err)
				continue
			}
			if len(payloadProperty.Values) == 0 {
				logrus.Errorf("empty values property: %s", payloadProperty.Id)
				continue
			}

			v := make([]struct {
				String  string  `json:"string,omitempty"`
				Integer int     `json:"integer,omitempty"`
				Float   float32 `json:"float,omitempty"`
				Boolean bool    `json:"boolean,omitempty"`
			}, len(payloadProperty.Values))
			for i, propertyValue := range payloadProperty.Values {
				parseValue, err := convertValueType(getProperty.Type, propertyValue)
				if err != nil {
					continue
				}
				v[i] = parseValue
			}
			createProperties[getProperty.Code] = struct {
				Id     string `json:"id"`
				Values []struct {
					String  string  `json:"string,omitempty"`
					Integer int     `json:"integer,omitempty"`
					Float   float32 `json:"float,omitempty"`
					Boolean bool    `json:"boolean,omitempty"`
				} `json:"values"`
			}{Id: getProperty.Id, Values: v}
		}

		if len(createProperties) > 0 {
			productUpdateFields.Properties = createProperties
		}
	}

	getProduct, err := p.GetById(productUpdateFields.Id)
	switch receivedMessage.Event {
	case "add":
		if err != nil {
			if err := p.Add(productUpdateFields); err != nil {
				return fmt.Errorf("could not add product: %w", err)
			}
		} else {
			return fmt.Errorf("could not add product becouse exist product %d", productUpdateFields.Id)
		}
	case "update":
		if err != nil {
			return fmt.Errorf("could not get product: %w", err)
		}
		if err := p.Update(getProduct.Id, productUpdateFields); err != nil {
			return fmt.Errorf("could not update product: %w", err)
		}
	case "delete":
		if err != nil {
			return fmt.Errorf("could not get product: %w", err)
		}
		if err := p.Delete(getProduct.Id); err != nil {
			return fmt.Errorf("could not delete product: %w", err)
		}
	default:
		return fmt.Errorf("unknown event: %s", receivedMessage.Event)
	}

	return nil
}

func convertValueType(propertyType string, propertyValue any) (valueStruct struct {
	String  string  `json:"string,omitempty"`
	Integer int     `json:"integer,omitempty"`
	Float   float32 `json:"float,omitempty"`
	Boolean bool    `json:"boolean,omitempty"`
}, err error) {
	switch propertyType {
	case "string":
		valueStruct.String = propertyValue.(string)
	case "integer":
		valueStruct.Integer = propertyValue.(int)
	case "float":
		valueStruct.Float = propertyValue.(float32)
	case "boolean":
		valueStruct.Boolean = propertyValue.(bool)
	default:
		return valueStruct, fmt.Errorf("error parse")
	}

	return valueStruct, nil
}
