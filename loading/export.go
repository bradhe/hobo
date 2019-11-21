package loading

import (
	"compress/gzip"
	"encoding/csv"
	"io"
	"net/url"
	"path"
	"strconv"

	"github.com/bradhe/location-search/content"
	"github.com/bradhe/location-search/models"
)

type ParseExportCallbackFunc func(models.City) error

func ParseExportReader(f io.Reader, cb ParseExportCallbackFunc) error {
	r := csv.NewReader(f)
	var headers []string

	for {
		li, err := r.Read()

		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if headers == nil {
			headers = li
		} else {
			c := models.City{
				ID:               getString("id", headers, li),
				Name:             getString("name", headers, li),
				ASCIIName:        getString("ascii_name", headers, li),
				Latitude:         getFloat64("latitude", headers, li),
				Longitude:        getFloat64("longitude", headers, li),
				Population:       getInt("population", headers, li),
				Timezone:         getString("timezone", headers, li),
				RegionID:         getString("region_id", headers, li),
				RegionName:       getString("region_name", headers, li),
				RegionASCIIName:  getString("region_ascii_name", headers, li),
				CountryID:        getString("country_id", headers, li),
				CountryName:      getString("country_name", headers, li),
				CountryASCIIName: getString("country_ascii_name", headers, li),
			}

			if err := cb(c); err != nil {
				break
			}
		}
	}

	return nil

}

func ParseExport(loc *url.URL, cb ParseExportCallbackFunc) error {
	f, err := content.NewReader(loc)

	if err != nil {
		return err
	}

	defer f.Close()

	// We separate control paths here because we have to explicitly manage
	// closing the underlying streams ourselves in the case that we're using
	// gzip.
	if path.Ext(loc.Path) == ".gz" {
		gz, err := gzip.NewReader(f)

		if err != nil {
			return err
		}

		defer gz.Close()

		return ParseExportReader(gz, cb)
	} else {
		return ParseExportReader(f, cb)
	}
}

func getString(header string, headers, row []string) string {
	for i, h := range headers {
		if header == h {
			return row[i]
		}
	}

	return ""
}

func parseFloat64(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

func getFloat64(header string, headers, row []string) float64 {
	return parseFloat64(getString(header, headers, row))
}

func parseInt(str string) int {
	i, _ := strconv.Atoi(str)
	return int(i)
}

func getInt(header string, headers, row []string) int {
	return parseInt(getString(header, headers, row))
}
