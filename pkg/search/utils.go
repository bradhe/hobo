package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

var schemeexp = regexp.MustCompile(`[a-zA-Z]://`)

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

func newBody(str string) io.ReadSeeker {
	return bytes.NewReader([]byte(str))
}

func query(str string) string {
	if b, err := json.Marshal(str); err != nil {
		panic(err)
	} else {
		return fmt.Sprintf(tmpl, string(b))
	}
}

func pickRandom(strs []string) string {
	// TODO: I think this hides a bug.
	if len(strs) == 0 {
		return ""
	}

	if len(strs) == 1 {
		return strs[0]
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	idx := r.Int() % len(strs)
	return strs[idx]
}

func addScheme(hosts []string) string {
	host := pickRandom(hosts)

	if schemeexp.MatchString(host) {
		return host
	}

	return "http://" + host
}

func addSchemes(hosts []string) []string {
	clean := make([]string, len(hosts))

	for i := range hosts {
		if schemeexp.MatchString(hosts[i]) {
			clean[i] = hosts[i]
		} else {
			clean[i] = "http://" + hosts[i]
		}
	}

	return clean
}

func dumpBody(resp *http.Response) string {
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)
	return string(buf)
}

type BackoffFunc func() error

func isInErrors(errs []error, err error) bool {
	for _, e := range errs {
		if e == err {
			return true
		}
	}

	return false
}

func sleep(i int) {
	mult := time.Duration((math.Exp(float64(i))/3.0)*1000.0) * time.Millisecond
	jit := time.Duration(rand.Intn(500)) * time.Millisecond

	logger.Debugf(" ... backing off %d", mult+jit)
	time.Sleep((mult + jit))
}

func doBackoff(fn BackoffFunc, errs ...error) error {
	var acc int

	for {
		acc += 1

		if err := fn(); err != nil {
			if isInErrors(errs, err) {
				// we wait and we keep going
				sleep(acc)
			} else {
				return err
			}
		} else {
			// err was nil if we get here.
			return nil
		}
	}
}
