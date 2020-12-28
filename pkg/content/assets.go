//go:generate go-bindata -o assets_data.go -pkg content ./assets/...
package content

import (
	"bytes"
	"io"
	"io/ioutil"
)

func AssetReader(name string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBuffer(MustAsset(name)))
}
