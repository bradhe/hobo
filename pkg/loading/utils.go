package loading

import "io/ioutil"
import "github.com/elastic/go-elasticsearch/esapi"
import "regexp"

var schemeexp = regexp.MustCompile(`[a-zA-Z]://`)

func dumpBody(resp *esapi.Response) string {
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)
	return string(buf)
}

func addSchemes(hosts []string) []string {
	clean := make([]string, 0, len(hosts))

	for i := range hosts {
		if schemeexp.MatchString(hosts[i]) {
			clean[i] = hosts[i]
		} else {
			clean[i] = "http://" + hosts[i]
		}
	}

	return clean
}
