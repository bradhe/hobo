package config

import "github.com/spf13/pflag"
import "github.com/spf13/viper"
import "github.com/joho/godotenv"

const DefaultAWSRegion = "us-west-2"

const DefaultAddr = "localhost:8080"

const DefaultUseEnv = false

const DefaultDebug = false

const DefaultAWSProfile = "hobo"

const DefaultElasticsearchHost = "localhost:9200"

const DefaultProfileIndexName = "profiles"

type AWS struct {
	// UseEnv indicates that we should load credentials from the EC2 role.
	UseEnv bool

	// Profile is the name of the shared profile to load.
	Profile string

	// We can supply credentials directly using this thing.
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

type Elasticsearch struct {
	Host          []string
	CityIndexName string
}

type Config struct {
	Args []string

	// Debug indicates if we want verbose log output and other debug-related
	// features enabled.
	Debug bool

	// Addr is the address to bind to when starting up the server.
	Addr string

	// AWS is the current AWS credential configuration.
	AWS AWS

	// Elasticsearch contains all the configuration related to elasticsearch
	Elasticsearch Elasticsearch

	// ExportURL is the URL that exports can be found (optional).
	ExportURL string

	// DataURL is the URL that data for parsing can be found (optional).
	DataURL string
}

func New() *Config {
	// NOTE: We do this first because we want the environment variables
	if err := godotenv.Load(); err != nil {
		logger.WithError(err).Warn("failed to load .env file")
	}

	viper.SetEnvPrefix("hobo")
	viper.AutomaticEnv()
	viper.SetDefault("debug", DefaultDebug)
	viper.SetDefault("addr", DefaultAddr)
	viper.SetDefault("aws_use_env", DefaultUseEnv)
	viper.SetDefault("aws_region", DefaultAWSRegion)
	viper.SetDefault("aws_profile", DefaultAWSProfile)
	viper.SetDefault("elasticsearch_host", DefaultElasticsearchHost)
	viper.SetDefault("elasticsearch_profile_index", DefaultProfileIndexName)

	pflag.Bool("debug", false, "puts the server in debug mode")

	pflag.Parse()

	viper.BindPFlag("debug", pflag.Lookup("debug"))

	return &Config{
		Args:      pflag.Args(),
		Debug:     viper.GetBool("debug"),
		Addr:      viper.GetString("addr"),
		ExportURL: viper.GetString("export-url"),
		DataURL:   viper.GetString("data-url"),
		AWS: AWS{
			Profile:         viper.GetString("aws_profile"),
			UseEnv:          viper.GetBool("aws_use_env"),
			AccessKeyID:     viper.GetString("aws_access_key_id"),
			SecretAccessKey: viper.GetString("aws_secret_access_key"),
			Region:          viper.GetString("aws_region"),
		},
		Elasticsearch: Elasticsearch{
			Host:          viper.GetStringSlice("elasticsearch_host"),
			CityIndexName: viper.GetString("elasticsearch_city_index"),
		},
	}
}
