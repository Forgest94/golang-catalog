package elasticSearch

import "github.com/elastic/go-elasticsearch/v8"

type Client struct {
	client *elasticsearch.TypedClient
}
