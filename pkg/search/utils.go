package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

var tmpl = `
{
  "sort": [
      { "_score": "desc" },
      { "Population": "desc" }
  ],
  "query": {
      "multi_match" : {
          "query" : %s,
          "type": "bool_prefix",
          "fields": ["Name^2", "ASCIIName^2", "RegionName", "RegionID", "CountryName", "CountryID"]
      }
  }
}`

func newBody(str string) io.Reader {
	return bytes.NewBufferString(str)
}

func query(str string) string {
	if b, err := json.Marshal(str); err != nil {
		panic(err)
	} else {
		return fmt.Sprintf(tmpl, string(b))
	}
}
