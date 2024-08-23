package store

import (
	"Golang/Go_Day03-1/src/types"
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Store interface {
	GetPlaces(limit, offset int) ([]types.Place, int, error)
}
type ElasticsearchStore struct {
	es        *elasticsearch.Client
	indexName string
}

func (store *ElasticsearchStore) NewElasticsearchStore(esClient *elasticsearch.Client, index string) *ElasticsearchStore {
	return &ElasticsearchStore{
		es:        esClient,
		indexName: index,
	}
}

type SearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source types.Place `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (store ElasticsearchStore) GetPlaces(limit, offset int) ([]types.Place, int, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, 0, err
	}

	searchRequest := esapi.SearchRequest{
		Index: []string{store.indexName},
		Body:  bytes.NewBuffer(queryBytes),
		Size:  &limit,
		From:  &offset,
		Sort:  []string{"id:asc"},
	}

	searchResult, err := searchRequest.Do(context.Background(), store.es)
	if err != nil {
		return nil, 0, err
	}
	defer searchResult.Body.Close()

	if searchResult.IsError() {
		return nil, 0, fmt.Errorf("search error: %s", searchResult.String())
	}

	var response SearchResponse

	if err := json.NewDecoder(searchResult.Body).Decode(&response); err != nil {
		return nil, 0, err
	}

	totalHits := response.Hits.Total.Value

	places := make([]types.Place, len(response.Hits.Hits))
	for i, hit := range response.Hits.Hits {
		places[i] = hit.Source
	}
	return places, totalHits, nil
}
