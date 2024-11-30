package property

import (
	"catalog/internal/kafka"
	"catalog/internal/services/property"
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
	topic = "properties"
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
		msg, _ := json.Marshal(property.ConsumerMessageV1{
			Uuid:    uuid.New().String(),
			Subject: "bd.property.service",
			Event:   events[rand.IntN(len(events))],
			Version: "1.0.0",
			Payload: struct {
				Id                string `json:"id"`
				IsActive          bool   `json:"is_active"`
				Name              string `json:"name"`
				Code              string `json:"code"`
				Hint              string `json:"hint,omitempty"`
				Type              string `json:"type"`
				ShowFilter        bool   `json:"show_filter"`
				ShowProductList   bool   `json:"show_product_list"`
				ShowDetailProduct bool   `json:"show_detail_product"`
			}{Id: productId, IsActive: true, Name: "property " + string(productId), Code: "property-" + string(productId), Hint: "test", Type: "string", ShowFilter: false, ShowProductList: false, ShowDetailProduct: false},
		})
		uuidKey := uuid.NewString()
		if err := p.Produce(topic, uuidKey, string(msg)); err != nil {
			logrus.Error(err)
		}
	}
}
