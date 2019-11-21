package parsing

import (
	"bufio"
	"io"
	"net/url"
	"strings"

	"github.com/bradhe/hobo/content"
	"github.com/bradhe/hobo/models"

	"github.com/mmcloughlin/geohash"
)

type Record struct {
	GeonameId             int
	Name                  string
	ASCIIName             string
	AlternateNames        []string
	Latitude              float64
	Longitude             float64
	FeatureClass          string
	FeatureCode           string
	CountryCode           string
	CC2                   string
	Admin1                string
	Admin2                string
	Admin3                string
	Admin4                string
	Population            int
	Elevation             int
	DigitalElevationModel string
	Timezone              string
	ModificationDate      string
}

func ParseRecord(row []string, r *Record) error {
	if len(row) < 19 {
		panic("invalid row length")
	}

	r.GeonameId = ParseInt(row[0])
	r.Name = row[1]
	r.ASCIIName = row[2]
	r.AlternateNames = strings.Split(row[3], ",")
	r.Latitude = ParseFloat64(row[4])
	r.Longitude = ParseFloat64(row[5])
	r.FeatureClass = row[6]
	r.FeatureCode = row[7]
	r.CountryCode = row[8]
	r.CC2 = row[9]
	r.Admin1 = row[10]
	r.Admin2 = row[11]
	r.Admin3 = row[12]
	r.Admin4 = row[13]
	r.Population = ParseInt(row[14])
	r.Elevation = ParseInt(row[15])
	r.DigitalElevationModel = row[16]
	r.Timezone = row[17]
	r.ModificationDate = row[18]
	return nil
}

func doParse(r io.Reader, fn func(*Record)) {
	var rec Record

	buf := bufio.NewReader(r)

	for {
		li, err := readLine(buf)

		if err != nil {
			if err != io.EOF {
				panic(err)
			}

			break
		}

		if err := ParseRecord(strings.Split(li, "\t"), &rec); err != nil {
			panic(err)
		} else {
			fn(&rec)
		}
	}
}

func LookupRegion(rec *Record, codes map[string]AdminCode) (string, string, string) {
	code, ok := codes[rec.Admin1]

	if !ok {
		code, ok = codes[rec.CountryCode+"."+rec.Admin1]

		if !ok {
			return "", "", ""
		}
	}

	return rec.Admin1, code.Name, code.ASCIIName
}

func Parse(loc *url.URL, cb func(*models.City)) error {
	countries := GetCountries()
	adminCodes := GetAdminCodes()

	r, err := content.NewReader(loc)

	if err != nil {
		return err
	}

	defer r.Close()

	doParse(r, func(rec *Record) {
		if rec.FeatureClass == "P" {
			regionID, regionName, regionASCIIName := LookupRegion(rec, adminCodes)

			city := models.City{
				ID:               geohash.Encode(rec.Latitude, rec.Longitude),
				Name:             rec.Name,
				ASCIIName:        rec.ASCIIName,
				Latitude:         rec.Latitude,
				Longitude:        rec.Longitude,
				Population:       rec.Population,
				Timezone:         rec.Timezone,
				RegionID:         regionID,
				RegionName:       regionName,
				RegionASCIIName:  regionASCIIName,
				CountryID:        rec.CountryCode,
				CountryName:      countries[rec.CountryCode].Name,
				CountryASCIIName: countries[rec.CountryCode].ASCIIName,
			}

			cb(&city)
		}
	})

	return nil
}
