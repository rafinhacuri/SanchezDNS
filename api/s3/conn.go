//nolint:gochecknoglobals
package s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

var client *s3.Client

func Connect() {
	endpoint := env.C.FsEndpoint

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("brazil"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				env.C.FsUser,
				env.C.FsPassword,
				"",
			),
		),
	)
	if err != nil {
		panic(err)
	}

	client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
		o.ResponseChecksumValidation = aws.ResponseChecksumValidationUnset
	})
}

func Test(ctx context.Context) error {
	_, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return err
	}

	return nil
}
