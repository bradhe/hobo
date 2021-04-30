exec curl -s -H 'Content-Type: application/json' --data-binary '{
	"sort": [
			{ "_score": "desc" },
			{ "Population": "desc" }
	],
	"query": {
		"multi_match" : {
			"query" : "New York",
			"type": "bool_prefix",
			"fields": ["Name^2", "ASCIIName^2", "RegionName", "RegionID", "CountryName", "CountryID"]
		}
	}
}' http://localhost:9200/cities/_search | json | less
