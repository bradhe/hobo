package loading

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/bradhe/location-search/models"
)

type BulkIndexBuffer struct {
	index string
	buf   *bytes.Buffer
}

func Bytes(r io.Reader) []byte {
	b, err := ioutil.ReadAll(r)

	if err != nil {
		panic(err)
	}

	return b
}

func (b BulkIndexBuffer) Add(c models.City) error {
	b.buf.WriteString(fmt.Sprintf(`{"index": {"_index": "%s", "_id": "%s"}}`, b.index, c.ID))
	b.buf.Write([]byte("\n"))
	b.buf.Write(Bytes(c.SerializeJSON()))
	return nil
}

func (b BulkIndexBuffer) Reader() io.Reader {
	return b.buf
}

func NewBulkIndexBuffer(index string) *BulkIndexBuffer {
	return &BulkIndexBuffer{index, bytes.NewBufferString("")}
}
