package search

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/bradhe/hobo/pkg/awsutils"
	"github.com/bradhe/hobo/pkg/config"
	"github.com/bradhe/hobo/pkg/models"

	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

type Client struct {
	conf *config.Config

	// the HTTP client to use
	c *http.Client
}

func (s *Client) indexUrl(name string) string {
	host := addScheme(s.conf.Elasticsearch.Host)

	if strings.HasSuffix(host, "/") {
		return host + name + "/_search"
	} else {
		return host + "/" + name + "/_search"
	}
}

func (s *Client) Search(place string) ([]models.City, error) {
	body := newBody(query(place))
	req, _ := http.NewRequest("POST", s.indexUrl("cities"), body)
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

func (s *Client) doImport(buf *BulkIndexBuffer) error {
	host := addScheme(s.conf.Elasticsearch.Host)
	body := buf.Reader()
	req, _ := http.NewRequest(http.MethodPost, host+"/_bulk", body)
	req.Header.Set("Content-Type", "application/json")

	if s.conf.Elasticsearch.IsSignedAuthentication() {
		logger.Debug("signing elasticsearch bulk write request")
		signer := v4.NewSigner(awsutils.Credentials(s.conf))
		signer.Sign(req, body, "es", s.conf.AWS.Region, time.Now())
	} else {
		logger.Debug("using anonymous bulk write request")
	}

	res, err := s.c.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		// A 403 indicates that there was a rate limiting error, so we should call
		// that out explicitly.
		if res.StatusCode == http.StatusForbidden {
			logger.WithField("body", dumpBody(res)).Error("import rate limited")
			return ErrImportRateLimited
		} else {
			logger.WithField("body", dumpBody(res)).Error("bulk write to ElasticSearch failed")
			return ErrImportFailed
		}
	}

	return nil
}

func (s *Client) Import(buf *BulkIndexBuffer) error {
	return doBackoff(func() error {
		return s.doImport(buf)
	}, ErrImportRateLimited)
}

func New(conf *config.Config) *Client {
	return &Client{conf, http.DefaultClient}
}
