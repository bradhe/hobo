package loading

import (
	"fmt"
	"net/url"

	"github.com/elastic/go-elasticsearch"
)

type Importer struct {
	es *elasticsearch.Client
}

func (i *Importer) Import(buf *BulkIndexBuffer) error {
	resp, err := i.es.Bulk(buf.Reader())

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("import failed. expected status 200, got status %d", resp.StatusCode)
	}

	return nil
}

func NewImporter(esurl *url.URL) (*Importer, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			esurl.String(),
		},
	}

	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		return nil, err
	}

	return &Importer{client}, nil
}
