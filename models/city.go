package models

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
)

type City struct {
	ID               string
	Name             string
	ASCIIName        string
	Latitude         float64
	Longitude        float64
	Population       int
	Timezone         string
	RegionID         string
	RegionName       string
	RegionASCIIName  string
	CountryID        string
	CountryName      string
	CountryASCIIName string
}

func (c City) WriteCSV(w *csv.Writer) error {
	rec := []string{
		c.ID,
		c.Name,
		c.ASCIIName,
		fmt.Sprintf("%f", c.Latitude),
		fmt.Sprintf("%f", c.Longitude),
		fmt.Sprintf("%d", c.Population),
		c.Timezone,
		c.RegionID,
		c.RegionName,
		c.RegionASCIIName,
		c.CountryID,
		c.CountryName,
		c.CountryASCIIName,
	}

	return w.Write(rec)
}

func (c City) SerializeJSON() io.Reader {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(c); err != nil {
		panic(err)
	}

	return &buf
}

func (c City) ToJSON() string {
	if buf, err := json.Marshal(&c); err != nil {
		panic(err)
	} else {
		return string(buf)
	}
}
