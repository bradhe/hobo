package content

import (
	"io"
	"net/url"

	"github.com/bradhe/hobo/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func s3Open(conf *config.Config, loc *url.URL) (io.ReadCloser, error) {
	sess := newAWSSession(conf)
	svc := s3.New(sess)

	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(loc.Host),
		Key:    aws.String(loc.Path),
	})

	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
