package property

import (
	"catalog/internal/elasticSearch"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Model struct {
	el *elasticSearch.Client
}

type Property struct {
	systemId          string
	Id                string `json:"id"`
	IsActive          bool   `json:"is_active"`
	Name              string `json:"name"`
	Code              string `json:"code"`
	Hint              string `json:"hint,omitempty"`
	Type              string `json:"type"`
	ShowFilter        bool   `json:"show_filter"`
	ShowProductList   bool   `json:"show_product_list"`
	ShowDetailProduct bool   `json:"show_detail_product"`
}

type ParamsQuery struct {
	Size        int
	Query       types.Query
	Aggregtions map[string]types.Aggregations
}
