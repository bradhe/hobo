package content

import (
	"github.com/bradhe/hobo/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
)

func newAWSSession(conf *config.Config) *session.Session {
	if conf.AWS.UseEnv {
		logger.Debug("configuratin AWS session from environment")

		return session.Must(session.NewSession(defaults.Get().Config))
	} else if conf.AWS.Profile != "" {
		logger.Debugf("configuratin AWS session with AWS profile `%s`", conf.AWS.Profile)

		return session.Must(session.NewSessionWithOptions(session.Options{
			Profile:           conf.AWS.Profile,
			SharedConfigState: session.SharedConfigEnable,
		}))
	}

	// This is the default configuration: Static credentials.
	logger.Debug("configuratin AWS session with static credentials")
	return session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(conf.AWS.Region),
			Credentials: credentials.NewStaticCredentials(conf.AWS.AccessKeyID, conf.AWS.SecretAccessKey, ""),
		},
	}))
}
