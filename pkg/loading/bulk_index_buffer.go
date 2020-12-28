package loading

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	"github.com/bradhe/hobo/pkg/models"
)

type BulkIndexBuffer struct {
	mut   sync.Mutex
	index string
	buf   *bytes.Buffer
	count int
}

func Bytes(r io.Reader) []byte {
	b, err := ioutil.ReadAll(r)

	if err != nil {
		panic(err)
	}

	return b
}

func (b *BulkIndexBuffer) Reset() {
	b.mut.Lock()
	defer b.mut.Unlock()

	b.count = 0
	b.buf.Reset()
}

func (b *BulkIndexBuffer) Add(c models.City) error {
	b.mut.Lock()
	defer b.mut.Unlock()

	b.buf.WriteString(fmt.Sprintf(`{"index": {"_index": "%s", "_id": "%s"}}`, b.index, c.ID))
	b.buf.Write([]byte("\n"))
	b.buf.WriteString(c.ToJSON())
	b.buf.Write([]byte("\n"))
	b.count += 1

	return nil
}

func (b *BulkIndexBuffer) Count() int {
	b.mut.Lock()
	defer b.mut.Unlock()

	return b.count
}

func (b *BulkIndexBuffer) Reader() io.Reader {
	return b.buf
}

func NewBulkIndexBuffer(index string) *BulkIndexBuffer {
	return &BulkIndexBuffer{sync.Mutex{}, index, bytes.NewBufferString(""), 0}
}
