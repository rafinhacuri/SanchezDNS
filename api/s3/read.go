package s3

import (
	"context"
	"errors"
	"io"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func Read(ctx context.Context, objectName string) ([]byte, string, error) {
	validKey := regexp.MustCompile(`^[a-zA-Z0-9._/-]+$`)
	if !validKey.MatchString(objectName) {
		return nil, "", errors.New("invalid object name")
	}

	out, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(env.C.FsBucket),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return nil, "", err
	}

	defer func() {
		err := out.Body.Close()
		if err != nil {
			log.Printf("failed to close S3 object body: %v", err)
		}
	}()

	data, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, "", err
	}

	return data, *out.ContentType, nil
}
