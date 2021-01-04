package search

import "errors"

var (
	ErrRateLimited   = errors.New("search: rate limited")
	ErrRequestFailed = errors.New("search: request failed")
)
