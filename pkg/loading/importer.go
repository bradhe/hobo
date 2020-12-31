package loading

import (
	"github.com/elastic/go-elasticsearch"
)

type Importer struct {
	es *elasticsearch.Client
}

func (i *Importer) Import(buf *BulkIndexBuffer) error {
	return nil
}

func NewImporter(hosts []string) (*Importer, error) {
	cfg := elasticsearch.Config{
		Addresses: addSchemes(hosts),
	}

	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		return nil, err
	}

	return &Importer{client}, nil
}
