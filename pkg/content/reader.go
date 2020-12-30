package content

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/bradhe/hobo/pkg/config"
)

func NewReader(conf *config.Config, loc *url.URL) (io.ReadCloser, error) {
	switch loc.Scheme {
	case "file":
		return os.Open(loc.Path)
	case "s3":
		return s3Open(conf, loc)
	}

	return nil, fmt.Errorf("content: %s reader not supported", loc.Scheme)
}
