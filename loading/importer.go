package loading

import (
	"net/url"

	"github.com/elastic/go-elasticsearch"
)

type Importer struct {
	es *elasticsearch.Client
}

func (i *Importer) Import(buf *BulkIndexBuffer) error {
	_, err := i.es.Bulk(buf.Reader())
	return err
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
