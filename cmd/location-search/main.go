package main

import (
	"encoding/csv"
	"net/url"
	"os"

	"github.com/bradhe/location-search/loading"
	"github.com/bradhe/location-search/models"
	"github.com/bradhe/location-search/parsing"
	"github.com/bradhe/location-search/search"
	"github.com/bradhe/location-search/server"

	"github.com/urfave/cli"
)

var gitCommit string
var version string

func doServe(c *cli.Context) error {
	go killOnSignal(c)

	logger.Infof("starting up location-search v%s (%s)", version, gitCommit)
	logger.Infof(" --addr=%s", c.GlobalString("addr"))
	logger.Infof(" --elasticsearch-addr=%s", c.GlobalString("elasticsearch-url"))

	client := search.New(c.GlobalString("elasticsearch-url"))
	server := server.New(client)

	return server.ListenAndServe(c.GlobalString("addr"))
}

func doLoad(c *cli.Context) error {
	go killOnSignal(c)

	esurl, err := url.Parse(c.GlobalString("elasticsearch-url"))

	if err != nil {
		panic(err)
	}

	importer, err := loading.NewImporter(esurl)

	if err != nil {
		panic(err)
	}

	export := c.String("export-url")

	buf := loading.NewBulkIndexBuffer("cities")

	logger.Infof("starting load for %s", export)

	var total int

	cb := func(city models.City) error {
		buf.Add(city)

		if buf.Count()%10000 == 0 {
			total += buf.Count()
			logger.Infof(" ... importing %d cities (%d total)", buf.Count(), total)

			if err = importer.Import(buf); err != nil {
				panic(err)
			} else {
				buf = loading.NewBulkIndexBuffer("cities")
			}
		}

		return nil
	}

	if export != "" {
		loc, err := url.Parse(export)

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

func killOnSignal(c *cli.Context) {
	// noop
}

func doParse(c *cli.Context) error {
	go killOnSignal(c)

	data := c.String("data-url")

	loc, err := url.Parse(data)

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

func main() {
	app := &cli.App{
		Name:  "location-search",
		Usage: "perform and maintain search against a",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "addr",
				Value: "localhost:8081",
				Usage: "address to bind http server to",
			},
			&cli.StringFlag{
				Name:  "elasticsearch-url",
				Value: "http://localhost:9200",
				Usage: "elasticsearch url to use for loading and search",
			},
		},
		Commands: []cli.Command{
			{
				Name:   "serve",
				Usage:  "serve an HTTP interface for serving",
				Action: doServe,
			},
			{
				Name:   "load",
				Usage:  "load data in to ElasticSearch",
				Action: doLoad,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "export-url",
						Value: "",
						Usage: "url of the export to import",
					},
				},
			},
			{
				Name:   "parse",
				Usage:  "parse raw data",
				Action: doParse,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "data-url",
						Value: "",
						Usage: "url of the data to normalize",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
