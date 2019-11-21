package server

import (
	"encoding/json"
)

func Dump(v interface{}) []byte {
	if b, err := json.MarshalIndent(v, "", "  "); err != nil {
		logger.WithError(err).Printf("dumping object failed")
		return []byte("")
	} else {
		return b
	}
}

func DumpString(v interface{}) string {
	return string(Dump(v))
}
