package search

import "errors"

var (
	ErrImportRateLimited = errors.New("search: import rate limited")
	ErrImportFailed      = errors.New("search: import failed")
)
