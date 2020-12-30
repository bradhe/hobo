package search

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bradhe/hobo/pkg/models"
)

type Client struct {
	// the URL of where ElasticSearch lives
	hosts []string

	// the HTTP client to use
	c *http.Client
}

func (s *Client) indexUrl(name string) string {
	host := addScheme(s.hosts)

	if strings.HasSuffix(host, "/") {
		return host + name + "/_search"
	} else {
		return host + "/" + name + "/_search"
	}
}

func (s *Client) Search(place string) ([]models.City, error) {
	req, _ := http.NewRequest("POST", s.indexUrl("cities"), newBody(query(place)))
	req.Header.Set("Content-Type", "application/json")
	res, err := s.c.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var results SearchResult

	if err := json.NewDecoder(res.Body).Decode(&results); err != nil {
		return nil, err
	}

	var cities []models.City

	for _, hit := range results.Hits.Hits {
		cities = append(cities, hit.City)
	}

	return cities, nil
}

func New(hosts []string) *Client {
	return &Client{hosts, http.DefaultClient}
}