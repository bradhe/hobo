package search

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bradhe/hobo/pkg/awsutils"
	"github.com/bradhe/hobo/pkg/config"
	"github.com/bradhe/hobo/pkg/models"

	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

type elasticsearchClient struct {
	conf *config.Config

	// the HTTP elasticsearchClient to use
	c *http.Client
}

func (s *elasticsearchClient) getCityPath(op string) string {
	index := s.conf.Elasticsearch.CityIndexName
	host := addScheme(s.conf.Elasticsearch.Host)

	if strings.HasSuffix(host, "/") {
		return host + index + "/" + op
	} else {
		return host + "/" + index + "/" + op
	}
}

func (s *elasticsearchClient) getOpPath(op string) string {
	host := addScheme(s.conf.Elasticsearch.Host)

	if strings.HasSuffix(host, "/") {
		return host + op
	} else {
		return host + "/" + op
	}
}

func (s *elasticsearchClient) doRequest(method string, path string, body io.ReadSeeker) (*http.Response, error) {
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")

	if s.conf.Elasticsearch.IsSignedAuthentication() {
		logger.Debug("signing elasticsearch search request")
		signer := v4.NewSigner(awsutils.Credentials(s.conf))
		signer.Sign(req, body, "es", s.conf.AWS.Region, time.Now())
	}

	res, err := s.c.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		// A 403 indicates that there was a rate limiting error, so we should call
		// that out explicitly.
		if res.StatusCode == http.StatusForbidden {
			logger.WithField("body", dumpBody(res)).Error("import rate limited")
			return nil, ErrRateLimited
		} else {
			logger.WithField("body", dumpBody(res)).Error("bulk write to ElasticSearch failed")
			return nil, ErrRequestFailed
		}
	}

	return res, nil
}

func (s *elasticsearchClient) Search(place string) ([]models.City, error) {
	body := newBody(query(place))
	res, err := s.doRequest(http.MethodPost, s.getCityPath("_search"), body)

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

func (s *elasticsearchClient) doImport(buf *BulkIndexBuffer) error {
	body := buf.Reader()

	res, err := s.doRequest(http.MethodPost, s.getOpPath("_bolk"), body)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	return err
}

func (s *elasticsearchClient) Import(buf *BulkIndexBuffer) error {
	return doBackoff(func() error {
		return s.doImport(buf)
	}, ErrRateLimited)
}

func newElasticsearchClient(conf *config.Config) Search {
	return &elasticsearchClient{conf, http.DefaultClient}
}
