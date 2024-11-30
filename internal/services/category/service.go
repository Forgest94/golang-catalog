package category

import (
	"catalog/internal/models/category"
	"encoding/json"
	"fmt"
)

func GetById(id string) (*category.Category, error) {
	c, _ := category.NewModel()
	return c.GetById(id)
}

func LoadMessageConsumerV1(message []byte) error {
	var receivedMessage ConsumerMessageV1
	err := json.Unmarshal(message, &receivedMessage)
	if err != nil {
		return fmt.Errorf("could not unmarshal category message v1: %w", err)
	}

	if receivedMessage.Payload.Id == "" {
		return fmt.Errorf("invalid category message v1")
	}

	c, _ := category.NewModel()
	categoryUpdateFields := category.Category{
		Id:          receivedMessage.Payload.Id,
		ParentId:    receivedMessage.Payload.ParentId,
		IsActive:    receivedMessage.Payload.IsActive,
		Name:        receivedMessage.Payload.Name,
		Code:        receivedMessage.Payload.Code,
		Description: receivedMessage.Payload.Description,
		ImageUrl:    receivedMessage.Payload.Img,
	}

	getCategory, err := c.GetById(categoryUpdateFields.Id)
	switch receivedMessage.Event {
	case "add":
		if err != nil {
			if err := c.Add(categoryUpdateFields); err != nil {
				return fmt.Errorf("could not add category: %w", err)
			}
		} else {
			return fmt.Errorf("could not add category becouse exist category %d", categoryUpdateFields.Id)
		}
	case "update":
		if err != nil {
			return fmt.Errorf("could not get category: %w", err)
		}
		if err := c.Update(getCategory.Id, categoryUpdateFields); err != nil {
			return fmt.Errorf("could not update category: %w", err)
		}
	case "delete":
		if err != nil {
			return fmt.Errorf("could not get category: %w", err)
		}
		if err := c.Delete(getCategory.Id); err != nil {
			return fmt.Errorf("could not delete category: %w", err)
		}
	default:
		return fmt.Errorf("unknown event: %s", receivedMessage.Event)
	}

	return nil
}
