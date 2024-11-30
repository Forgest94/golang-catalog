package product

import (
	"catalog/internal/elasticSearch"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Model struct {
	el *elasticSearch.Client
}

type Product struct {
	systemId           string
	Id                 string   `json:"id"`
	ParentId           string   `json:"parent_id"`
	IsActive           bool     `json:"is_active"`
	Name               string   `json:"name"`
	Code               string   `json:"code"`
	Description        string   `json:"description"`
	PreviewDescription string   `json:"preview_description"`
	ImageUrl           string   `json:"image_url"`
	BasePrice          float32  `json:"base_price"`
	Price              float32  `json:"price"`
	Categories         []string `json:"categories,omitempty"`
	Properties         map[string]struct {
		Id     string `json:"id"`
		Values []struct {
			String  string  `json:"string,omitempty"`
			Integer int     `json:"integer,omitempty"`
			Float   float32 `json:"float,omitempty"`
			Boolean bool    `json:"boolean,omitempty"`
		} `json:"values"`
	} `json:"properties,omitempty"`
}

type ParamsQuery struct {
	Size        int
	Query       types.Query
	Aggregtions map[string]types.Aggregations
}
