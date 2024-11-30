package elasticSearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"io"
)

func NewClient(address []string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: address,
	}

	c, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("create elasticsearch client: %w", err)
	}

	return &Client{
		client: c,
	}, nil
}

func (c *Client) ExistIndex(index string) (bool, error) {
	res, err := c.client.Indices.Exists(index).Do(context.Background())
	if err != nil {
		return false, fmt.Errorf("check exist index: %w", err)
	}

	return res, err
}

func (c *Client) CreateIndex(index string, propertiesIndex map[string]types.Property) error {
	existIndex, err := c.ExistIndex(index)
	if err != nil {
		return fmt.Errorf("check exist index: %w", err)
	}

	if existIndex {
		return nil
	}

	_, err = c.client.Indices.Create(index).Request(&create.Request{
		Mappings: &types.TypeMapping{
			Properties: propertiesIndex,
		},
	}).Do(nil)

	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}

	return nil
}

func (c *Client) DeleteIndex(index string) error {
	_, err := c.client.Indices.Delete(index).Do(context.Background())
	if err != nil {
		return fmt.Errorf("delete index: %w", err)
	}

	return nil
}

func (c *Client) CreateDocument(index string, doc interface{}) error {
	strDoc, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("create error serialize: %w", err)
	}

	_, err = c.client.Index(index).Raw(io.Reader(bytes.NewBuffer(strDoc))).Do(context.Background())
	if err != nil {
		return fmt.Errorf("create document: %w", err)
	}

	return nil
}

func (c *Client) Search(index string, searchParams search.Request) (*search.Response, error) {
	res, err := c.client.
		Search().
		Index(index).
		Request(&searchParams).
		Do(context.Background())

	if err != nil {
		return res, fmt.Errorf("error searching documents: %w", err)
	}

	return res, nil
}

func (c *Client) UpdateDocument(index string, id string, doc interface{}) error {
	strDoc, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("update error serialize: %w", err)
	}

	_, err = c.client.Update(index, id).
		Request(&update.Request{Doc: strDoc}).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}

	return nil
}

func (c *Client) DeleteDocument(index string, id string) error {
	_, err := c.client.Delete(index, id).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting document: %w", err)
	}

	return nil
}
