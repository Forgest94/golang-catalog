package indexing

import (
	"catalog/internal/elasticSearch"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

const (
	indexProducts   = "products"
	indexCategories = "categories"
	indexProperties = "properties"
)

var (
	address         = strings.Split(os.Getenv("ELASTICSEARCH_HOSTS"), ",")
	productsMapping = map[string]types.Property{
		"id":                  types.NewKeywordProperty(),
		"parent_id":           types.NewKeywordProperty(),
		"is_active":           types.NewBooleanProperty(),
		"name":                types.NewKeywordProperty(),
		"code":                types.NewKeywordProperty(),
		"description":         types.NewTextProperty(),
		"preview_description": types.NewTextProperty(),
		"image_url":           types.NewTextProperty(),
		"base_price":          types.NewFloatNumberProperty(),
		"price":               types.NewFloatNumberProperty(),
		"categories":          types.NewKeywordProperty(),
		"properties":          types.NewObjectProperty(),
	}
	categoriesMapping = map[string]types.Property{
		"id":          types.NewKeywordProperty(),
		"parent_id":   types.NewKeywordProperty(),
		"is_active":   types.NewBooleanProperty(),
		"name":        types.NewKeywordProperty(),
		"code":        types.NewKeywordProperty(),
		"description": types.NewTextProperty(),
		"image_url":   types.NewTextProperty(),
	}
	propertiesMapping = map[string]types.Property{
		"id":                  types.NewKeywordProperty(),
		"is_active":           types.NewBooleanProperty(),
		"name":                types.NewKeywordProperty(),
		"code":                types.NewKeywordProperty(),
		"hint":                types.NewTextProperty(),
		"type":                types.NewKeywordProperty(),
		"show_filter":         types.NewBooleanProperty(),
		"show_product_list":   types.NewBooleanProperty(),
		"show_detail_product": types.NewBooleanProperty(),
	}
)

func Run(w http.ResponseWriter, r *http.Request) {
	e, err := elasticSearch.NewClient(address)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := e.CreateIndex(indexProducts, productsMapping); err != nil {
		logrus.Fatal(err)
	}

	if err := e.CreateIndex(indexCategories, categoriesMapping); err != nil {
		logrus.Fatal(err)
	}

	if err := e.CreateIndex(indexProperties, propertiesMapping); err != nil {
		logrus.Fatal(err)
	}
}
