package content

import (
	"fmt"
	"io"
	"net/url"
	"os"
)

func newS3Reader(loc *url.URL) (io.ReadCloser, error) {
	return nil, fmt.Errorf("content: s3 reader not implemented")
}

func NewReader(loc *url.URL) (io.ReadCloser, error) {
	switch loc.Scheme {
	case "file":
		return os.Open(loc.Path)
	case "s3":
		return newS3Reader(loc)
	}

	return nil, fmt.Errorf("content: %s reader not supported", loc.Scheme)
}
