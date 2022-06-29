package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
	"strings"
)

const (
	imageJPG = "image/jpeg"
	imagePNG = "image/png"
	fontTTF  = "font/ttf"
)

var (
	// ErrAccessKeyIDNotFound is returned when the AWS Access Key ID can't be
	// found in the process's environment.

	// ErrSecretAccessKeyNotFound is returned when the AWS Secret Access Key
	// can't be found in the process's environment.
	ErrSecretAccessKeyNotFound = awserr.New("EnvSecretOrAccessNotFound", "AWS_SECRET_ACCESS_KEY or AWS_ACCESS_KEY_ID not found in environment", nil)
)

const (
	ACCESS_ID  = "AWS_ACCESS_KEY_ID"
	SECRET_KEY = "AWS_SECRET_ACCESS_KEY"
)

type PutFile struct {
	Bucket          string
	FileName        string
	FileDestination string
	ContentType     string
	BucketOwner     string
	Body            io.Reader
	ContentLength   int
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Cant read env vars. %s \n", err.Error())
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "yc",
			URL:           "https://storage.yandexcloud.net",
			SigningRegion: "us-east-1",
		}, nil
	})

	envProvider := NewEnvCredentials()

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver), config.WithCredentialsProvider(envProvider))

	if err != nil {
		log.Fatal(err)
	}
	client := s3.NewFromConfig(cfg)

	bucketName := "zharpizza-bucket"

	path := "upload/upload_me.jpg"
	uploadPath := "static/filename123.jpg"
	if _, err = os.Stat(path); err != nil {
		if !os.IsExist(err) {
			log.Fatal(err.Error())
		}
	}

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}
	conttype := "image/jpeg"
	ownr := "aje4n7m8g43hrki88uf7"
	grant := ""
	output, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:              &bucketName,
		Key:                 &uploadPath,
		Body:                bytes.NewReader(fileBytes),
		ContentType:         &conttype,
		ExpectedBucketOwner: &ownr,
		GrantRead:           &grant,
		ContentLength:       int64(len(fileBytes)),
	})
	if err != nil {
		log.Fatal(err.Error(), "here")
	}
	fmt.Println(*output.ETag)

}

func getFileNameByExt(ext string, file string) (string, bool) {
	split := strings.Split(file, ext)
	splitBySlash := strings.Split(strings.Join(split, ""), "/")
	afterSlash := splitBySlash[len(splitBySlash)-1]
	if strings.TrimSpace(afterSlash) == "" {
		return "", false
	}

	if l := len(strings.Split(afterSlash, ".")); l > 1 {
		return "", false
	}
	return afterSlash + ext, true
}

type EnvProvider struct {
	retrieved bool
}

func NewEnvCredentials() *EnvProvider {
	return &EnvProvider{}
}

func (p *EnvProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {

	accessId := os.Getenv(ACCESS_ID)
	secretKey := os.Getenv(SECRET_KEY)
	if accessId == "" || secretKey == "" {
		return aws.Credentials{}, ErrSecretAccessKeyNotFound
	}

	p.retrieved = true

	return aws.Credentials{
		AccessKeyID:     accessId,
		SecretAccessKey: secretKey,
	}, nil

}

func (p *EnvProvider) IsExpired() bool {
	return !p.retrieved
}
