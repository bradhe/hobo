package loading

import (
	"fmt"

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
		logger.WithField("body", dumpBody(resp)).Error("bulk write to ElasticSearch failed")
		return fmt.Errorf("import failed. expected status 200, got status %d", resp.StatusCode)
	}

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
