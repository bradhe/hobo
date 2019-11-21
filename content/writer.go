package content

import (
	"fmt"
	"io"
	"net/url"
	"os"
)

func newS3Writer(loc *url.URL) (io.WriteCloser, error) {
	return nil, fmt.Errorf("content: s3 reader not implemented")
}

func NewWriter(loc *url.URL) (io.WriteCloser, error) {
	switch loc.Scheme {
	case "file":
		return os.Open(loc.Path)
	case "s3":
		return newS3Writer(loc)
	}

	return nil, fmt.Errorf("content: %s writer not supported", loc.Scheme)
}
