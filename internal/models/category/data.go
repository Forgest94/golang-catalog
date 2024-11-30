package category

import (
	"catalog/internal/elasticSearch"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Model struct {
	el *elasticSearch.Client
}

type Category struct {
	systemId    string
	Id          string `json:"id"`
	ParentId    string `json:"parent_id,omitempty"`
	IsActive    bool   `json:"is_active"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url,omitempty"`
}

type ParamsQuery struct {
	Size        int
	Query       types.Query
	Aggregtions map[string]types.Aggregations
}
