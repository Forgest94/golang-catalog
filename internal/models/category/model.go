package category

import (
	"catalog/internal/elasticSearch"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/some"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

const indexName = "categories"

var address = strings.Split(os.Getenv("ELASTICSEARCH_HOSTS"), ",")

func NewModel() (*Model, error) {
	el, err := elasticSearch.NewClient(address)
	if err != nil {
		logrus.Fatal(err)
	}
	return &Model{
		el,
	}, nil
}

func (m *Model) GetById(id string) (*Category, error) {
	request := search.Request{
		Size: some.Int(1),
		Query: &types.Query{
			Match: map[string]types.MatchQuery{
				"id": {Query: id},
			},
		},
	}

	res, err := m.el.Search(indexName, request)
	if err != nil {
		return nil, err
	}

	if res.Hits.Total.Value == 0 {
		return nil, errors.New("category not found")
	}

	var category Category
	for _, doc := range res.Hits.Hits {
		if err := json.Unmarshal(doc.Source_, &category); err != nil {
			return nil, err
		}
	}

	return &category, nil
}

func (m *Model) Add(categoryFields Category) error {
	updateCategory, err := json.Marshal(categoryFields)
	if err != nil {
		return err
	}

	var category Category
	if err := json.Unmarshal(updateCategory, &category); err != nil {
		return err
	}

	if err := m.el.CreateDocument(indexName, category); err != nil {
		return err
	}

	return nil
}

func (m *Model) Update(id string, categoryFields Category) error {
	exists, err := m.GetById(id)
	if err != nil {
		if err := m.Add(categoryFields); err != nil {
			return err
		}
	} else {
		updateCategory, err := json.Marshal(categoryFields)
		if err != nil {
			return err
		}

		var category Category
		if err := json.Unmarshal(updateCategory, &category); err != nil {
			return err
		}

		if err := m.el.UpdateDocument(indexName, exists.systemId, category); err != nil {
			return err
		}
	}

	return nil
}

func (m *Model) Delete(id string) error {
	exists, err := m.GetById(id)
	if err != nil {
		return err
	}

	if err := m.el.DeleteDocument(indexName, exists.systemId); err != nil {
		return err
	}

	return nil
}

func (m *Model) Get(params ParamsQuery) ([]Category, map[string]struct{ Buckets []map[string]any }, error) {
	responseSearch, err := m.el.Search(indexName, search.Request{
		Size:         some.Int(params.Size),
		Query:        &params.Query,
		Aggregations: params.Aggregtions,
	})

	if err != nil {
		return nil, nil, err
	}

	responseElements := responseSearch.Hits.Hits
	responseAggregations := responseSearch.Aggregations

	categories := make([]Category, len(responseElements))
	for _, hit := range responseElements {
		var category Category
		if err := json.Unmarshal(hit.Source_, &category); err != nil {
			continue
		}
		categories = append(categories, category)
	}

	resultAggs := make(map[string]struct {
		Buckets []map[string]any
	})
	for code, agg := range responseAggregations {
		var Bucket struct {
			Buckets []map[string]any
		}
		a, err := json.Marshal(agg)
		if err != nil {
			continue
		}

		if err := json.Unmarshal(a, &Bucket); err != nil {
			continue
		}
		resultAggs[code] = Bucket
	}

	return categories, resultAggs, nil
}
