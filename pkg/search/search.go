package search

import "github.com/bradhe/hobo/pkg/models"
import "github.com/bradhe/hobo/pkg/config"

type Search interface {
	Search(string) ([]models.City, error)
	Import(*BulkIndexBuffer) error
}

func New(conf *config.Config) Search {
	return newElasticsearchClient(conf)
}
