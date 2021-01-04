package awsutils

import (
	"github.com/bradhe/hobo/pkg/config"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
)

func Credentials(conf *config.Config) *credentials.Credentials {
	if conf.AWS.UseEnv {
		logger.Debug("fetching AWS credentials from environment")
		return defaults.Get().Config.Credentials
	} else if conf.AWS.Profile != "" {
		logger.Debugf("fetching AWS credentials from profile `%s`", conf.AWS.Profile)
		return credentials.NewSharedCredentials("", conf.AWS.Profile)
	}

	// This is the default configuration: Static credentials.
	logger.Debug("using static AWS credentials")
	return credentials.NewStaticCredentials(conf.AWS.AccessKeyID, conf.AWS.SecretAccessKey, "")
}
