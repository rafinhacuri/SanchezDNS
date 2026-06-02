package s3

import (
	"bufio"
	"context"
	"errors"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/text/unicode/norm"

	"github.com/rafinhacuri/SanchezDNS/api/env"
)

func getFileName(fileName string) string {
	decoded, _ := url.QueryUnescape(fileName)
	normalized := norm.NFC.String(decoded)

	return path.Base(normalized)
}

func createFileName(fileName string) string {
	baseName := strings.TrimSpace(getFileName(fileName))

	origExt := path.Ext(baseName)
	ext := strings.ToLower(origExt)
	name := strings.TrimSuffix(baseName, origExt)

	normalized := norm.NFD.String(name)

	reAccents := regexp.MustCompile("[\u0300-\u036F]")
	finalName := reAccents.ReplaceAllString(normalized, "")

	reInvalid := regexp.MustCompile(`[^\w-]`)
	finalName = reInvalid.ReplaceAllString(finalName, "-")

	reMultiDash := regexp.MustCompile(`-+`)
	finalName = reMultiDash.ReplaceAllString(finalName, "-")

	finalName = strings.ToLower(finalName)
	finalName = strings.Trim(finalName, "-")

	return finalName + "-" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ext
}

func Save(ctx context.Context, file *multipart.FileHeader, prefix string) (string, error) {
	if file == nil {
		return "", errors.New("file is nil")
	}

	src, err := file.Open()
	if err != nil {
		log.Println("error opening file:", err)

		return "", errors.New("erro interno")
	}

	defer func() {
		err = src.Close()
		if err != nil {
			log.Println("error closing file:", err)
		}
	}()

	contentType := strings.TrimSpace(file.Header.Get("Content-Type"))

	br := bufio.NewReader(src)
	if contentType == "" || contentType == "application/octet-stream" {
		head, _ := br.Peek(512)
		if len(head) > 0 {
			contentType = http.DetectContentType(head)
		}
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	objectName := createFileName(file.Filename)
	if prefix = strings.Trim(prefix, "/"); prefix != "" {
		objectName = prefix + "/" + objectName
	}

	size := max(file.Size, 0)

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(env.C.FsBucket),
		Key:           aws.String(objectName),
		Body:          br,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(size),
	})
	if err != nil {
		log.Println("error putting object to s3:", err)

		return "", errors.New("erro interno")
	}

	return objectName, nil
}
