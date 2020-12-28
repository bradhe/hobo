package loading

import "io/ioutil"
import "github.com/elastic/go-elasticsearch/esapi"

func dumpBody(resp *esapi.Response) string {
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)
	return string(buf)
}
