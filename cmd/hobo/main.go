package main

import (
	"encoding/csv"
	"net/url"
	"os"
	"strings"

	"github.com/bradhe/hobo/pkg/config"
	"github.com/bradhe/hobo/pkg/loading"
	"github.com/bradhe/hobo/pkg/logs"
	"github.com/bradhe/hobo/pkg/models"
	"github.com/bradhe/hobo/pkg/parsing"
	"github.com/bradhe/hobo/pkg/search"
	"github.com/bradhe/hobo/pkg/server"
)

var gitCommit string
var version string

func doPrintStartup(command string, conf *config.Config) {
	logger.Infof("starting up hobo v%s (%s)", version, gitCommit)
	logger.Infof(" --debug=%b", conf.Debug)
	logger.Infof(" --addr=%s", conf.Addr)
	logger.Infof(" --elasticsearch-addr=%s", strings.Join(conf.Elasticsearch.Host, ","))
	logger.Infof(" --elasticsearch-authentication=%s", conf.Elasticsearch.Authentication)
	logger.Infof("executing: %s", command)
}

func doServe(conf *config.Config) error {
	return server.New(conf).ListenAndServe(conf.Addr)
}

func doImport(conf *config.Config) error {
	client := search.New(conf)
	buf := search.NewBulkIndexBuffer(conf.Elasticsearch.CityIndexName)

	logger.Infof("starting load for %s", conf.ExportURL)

	var total int

	op := func(city models.City) error {
		buf.Add(city)

		if buf.Count()%1000 == 0 {
			total += buf.Count()
			logger.Infof(" ... importing %d cities (%d total)", buf.Count(), total)

			if err := client.Import(buf); err != nil {
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

		if err := loading.ParseExport(conf, loc, op); err != nil {
			panic(err)
		}
	} else {
		if err := loading.ParseExportReader(os.Stdin, op); err != nil {
			panic(err)
		}
	}

	logger.Infof("imported complete. imported %d cities", total)

	return nil
}

func doParse(conf *config.Config) error {
	loc, err := url.Parse(conf.DataURL)

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

	parsing.Parse(conf, loc, func(city *models.City) {
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

	if conf.Debug {
		logs.EnableDebug()
	}

	cmd := GetCommand(conf)

	doPrintStartup(cmd, conf)

	switch cmd {
	case "import":
		doImport(conf)
	case "serve":
		fallthrough
	default:
		doServe(conf)
	}
}
