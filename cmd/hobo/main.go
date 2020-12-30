package main

import (
	"encoding/csv"
	"net/url"
	"os"
	"strings"

	"github.com/bradhe/hobo/pkg/config"
	"github.com/bradhe/hobo/pkg/loading"
	"github.com/bradhe/hobo/pkg/models"
	"github.com/bradhe/hobo/pkg/parsing"
	"github.com/bradhe/hobo/pkg/search"
	"github.com/bradhe/hobo/pkg/server"
)

var gitCommit string
var version string

func doServe(conf *config.Config) error {
	logger.Infof("starting up hobo v%s (%s)", version, gitCommit)
	logger.Infof(" --addr=%s", conf.Addr)
	logger.Infof(" --elasticsearch-addr=%s", strings.Join(conf.Elasticsearch.Host, ","))

	client := search.New(conf.Elasticsearch.Host)
	server := server.New(client)

	return server.ListenAndServe(conf.Addr)
}

func doImport(conf *config.Config) error {
	importer, err := loading.NewImporter(conf.Elasticsearch.Host)

	if err != nil {
		panic(err)
	}

	buf := loading.NewBulkIndexBuffer("cities")

	logger.Infof("starting load for %s", conf.ExportURL)

	var total int

	cb := func(city models.City) error {
		buf.Add(city)

		if buf.Count()%5000 == 0 {
			total += buf.Count()
			logger.Infof(" ... importing %d cities (%d total)", buf.Count(), total)

			if err = importer.Import(buf); err != nil {
				panic(err)
			} else {
				buf.Reset()
			}
		}

		return nil
	}

	if conf.ExportURL != "" {
		loc, err := url.Parse(conf.ExportURL)

		if err != nil {
			panic(err)
		}

		if err := loading.ParseExport(loc, cb); err != nil {
			panic(err)
		}
	} else {
		if err := loading.ParseExportReader(os.Stdin, cb); err != nil {
			panic(err)
		}
	}

	return nil
}

func doParse(c *config.Config) error {
	loc, err := url.Parse(c.DataURL)

	if err != nil {
		panic(err)
	}

	output := csv.NewWriter(os.Stdout)
	output.Write([]string{
		"id",
		"name",
		"ascii_name",
		"latitude",
		"longitude",
		"population",
		"timezone",
		"region_id",
		"region_name",
		"region_ascii_name",
		"country_id",
		"country_name",
		"country_ascii_name",
	})

	defer output.Flush()

	parsing.Parse(loc, func(city *models.City) {
		city.WriteCSV(output)

	})

	return nil
}

func GetCommand(conf *config.Config) string {
	if len(conf.Args) > 0 {
		return conf.Args[0]
	}

	return ""
}

func main() {
	conf := config.New()

	cmd := GetCommand(conf)

	switch cmd {
	case "import":
		doImport(conf)
	case "serve":
		fallthrough
	default:
		doServe(conf)
	}
}
