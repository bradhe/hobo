package search

import (
	"github.com/bradhe/hobo/pkg/models"
)

type SearchHit struct {
	Index string      `json:"_index"`
	Id    string      `json:"_id"`
	Type  string      `json:"_type"`
	City  models.City `json:"_source"`
}

type SearchHits struct {
	Total struct {
		Value int `json:"value"`
	} `json:"total"`
	MaxScore *int        `json:"max_score"`
	Hits     []SearchHit `json:"hits"`
}

type SearchResult struct {
	TimedOut bool       `json:"timed_out"`
	Hits     SearchHits `json:"hits"`
}
