package category

import (
	"catalog/internal/kafka"
	"catalog/internal/services/category"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	topic = "categories"
	group
)

var address = strings.Split(os.Getenv("KAFKA_HOSTS"), ",")

func SendMessages(w http.ResponseWriter, r *http.Request) {
	p, err := kafka.NewProducer(address, group)
	if err != nil {
		logrus.Fatal(err)
	}

	events := []string{"add", "update", "delete"}

	for i := 0; i < 100; i++ {
		categoryId := strings.Join([]string{uuid.New().String(), strconv.Itoa(i)}, "-")
		msg, _ := json.Marshal(category.ConsumerMessageV1{
			Uuid:    uuid.New().String(),
			Subject: "bd.category.service",
			Event:   events[rand.IntN(len(events))],
			Version: "1.0.0",
			Payload: struct {
				Id          string `json:"id"`
				ParentId    string `json:"parent_id,omitempty"`
				IsActive    bool   `json:"is_active"`
				Name        string `json:"name"`
				Code        string `json:"code"`
				Description string `json:"description"`
				Img         string `json:"img,omitempty"`
			}{Id: categoryId, ParentId: uuid.New().String(), IsActive: true, Name: "cateory " + string(categoryId), Code: "category-" + string(categoryId), Description: "test", Img: "logo/img.png"},
		})
		uuidKey := uuid.NewString()
		if err := p.Produce(topic, uuidKey, string(msg)); err != nil {
			logrus.Error(err)
		}
	}
}
