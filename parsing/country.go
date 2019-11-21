package parsing

import (
	"encoding/json"

	"github.com/bradhe/location-search/content"
)

type Country struct {
	Name         string `json:"name"`
	ASCIIName    string `json:"ascii_name"`
	OfficialName string `json:"official_name"`
	Alpha2       string `json:"alpha2"`
	Alpha3       string `json:"alpha3"`
}

func GetCountries() map[string]Country {
	var countries []Country

	dec := json.NewDecoder(content.AssetReader("assets/countries.json"))

	if err := dec.Decode(&countries); err != nil {
		panic(err)
	}

	lookup := make(map[string]Country, len(countries))

	for _, country := range countries {
		lookup[country.Alpha2] = country
	}

	return lookup
}
