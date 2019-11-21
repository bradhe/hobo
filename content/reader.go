package content

import (
	"fmt"
	"io"
	"net/url"
	"os"
)

func NewReader(loc *url.URL) (io.ReadCloser, error) {
	switch loc.Scheme {
	case "file":
		return os.Open(loc.Path)
	case "s3":
		return s3Open(loc)
	}

	return nil, fmt.Errorf("content: %s reader not supported", loc.Scheme)
}
