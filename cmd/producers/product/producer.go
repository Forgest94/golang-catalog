package product

import (
	"catalog/internal/kafka"
	"catalog/internal/services/product"
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
	topic = "products"
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
		productId := strings.Join([]string{uuid.New().String(), strconv.Itoa(i)}, "-")
		msg, _ := json.Marshal(product.ConsumerMessageV1{
			Uuid:    uuid.New().String(),
			Subject: "bd.product.service",
			Event:   events[rand.IntN(len(events))],
			Version: "1.0.0",
			Payload: struct {
				Id                 string   `json:"id"`
				ParentId           string   `json:"parent_id,omitempty"`
				IsActive           bool     `json:"is_active"`
				Name               string   `json:"name"`
				Code               string   `json:"code"`
				Description        string   `json:"description"`
				PreviewDescription string   `json:"preview_description"`
				Img                string   `json:"img"`
				Price              float32  `json:"price"`
				BasePrice          float32  `json:"base_price"`
				Categories         []string `json:"categories,omitempty"`
				Properties         []struct {
					Id     string `json:"id"`
					Values []any  `json:"values"`
				} `json:"properties,omitempty"`
			}{Id: productId, ParentId: uuid.New().String(), IsActive: true, Name: "product " + string(productId), Code: "product-" + string(productId), Description: "fullText", PreviewDescription: "smallText", Img: "logo/img.png", Price: rand.Float32(), BasePrice: rand.Float32()},
		})
		uuidKey := uuid.NewString()
		if err := p.Produce(topic, uuidKey, string(msg)); err != nil {
			logrus.Error(err)
		}
	}
}
