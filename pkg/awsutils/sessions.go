package awsutils

import (
	"github.com/bradhe/hobo/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func Session(conf *config.Config) *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(conf.AWS.Region),
			Credentials: Credentials(conf),
		},
	}))
}
