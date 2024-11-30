package property

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

const indexName = "properties"

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

func (m *Model) GetById(id string) (*Property, error) {
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
		return nil, errors.New("property not found")
	}

	var property Property
	for _, doc := range res.Hits.Hits {
		if err := json.Unmarshal(doc.Source_, &property); err != nil {
			return nil, err
		}
	}

	return &property, nil
}

func (m *Model) Add(propertyFields Property) error {
	updateProperty, err := json.Marshal(propertyFields)
	if err != nil {
		return err
	}

	var property Property
	if err := json.Unmarshal(updateProperty, &property); err != nil {
		return err
	}

	if err := m.el.CreateDocument(indexName, property); err != nil {
		return err
	}

	return nil
}

func (m *Model) Update(id string, propertyFields Property) error {
	exists, err := m.GetById(id)
	if err != nil {
		if err := m.Add(propertyFields); err != nil {
			return err
		}
	} else {
		updateProperty, err := json.Marshal(propertyFields)
		if err != nil {
			return err
		}

		var property Property
		if err := json.Unmarshal(updateProperty, &property); err != nil {
			return err
		}

		if err := m.el.UpdateDocument(indexName, exists.systemId, property); err != nil {
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

func (m *Model) Get(params ParamsQuery) ([]Property, map[string]struct{ Buckets []map[string]any }, error) {
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

	properties := make([]Property, len(responseElements))
	for _, hit := range responseElements {
		var property Property
		if err := json.Unmarshal(hit.Source_, &property); err != nil {
			continue
		}
		properties = append(properties, property)
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

	return properties, resultAggs, nil
}
