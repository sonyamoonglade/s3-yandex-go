package s3yandex

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
)

type YandexS3Client struct {
	client *s3.Client
	bucket string
	owner  string
	debug  bool
	logger *log.Logger
}

func (y *YandexS3Client) PutFile(ctx context.Context, inp *PutFileInput) error {

	fileBytes, err := GetFileBytes(inp.FilePath, inp.FileName)
	if err != nil {
		return err
	}

	_, err = y.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:              &y.bucket,
		Key:                 &inp.Destination,
		Body:                bytes.NewReader(fileBytes),
		ContentType:         &inp.ContentType,
		ExpectedBucketOwner: &y.owner,
		GrantRead:           &grant,
		ContentLength:       int64(len(fileBytes)),
	})

	if err != nil {
		return err
	}
	y.logger.Println(fmt.Sprintf("[YANDEX S3] put file from '%s%s' into-> '%s'. bucket - %s", inp.FilePath, inp.FileName, inp.Destination, y.bucket))
	return nil
}

func NewYandexS3Client(provider aws.CredentialsProvider, yndxConfig YandexS3Config) *YandexS3Client {
	logger := log.New(os.Stdout, "[DEBUG] ", 0)
	// setting custom resolver that is specific for yandexcloud
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "yc",
			URL:           "https://storage.yandexcloud.net",
			SigningRegion: "us-east-1",
		}, nil
	})
	if yndxConfig.Debug {
		logger.Println("[YANDEX S3] initialized resolver")
	}
	//loading config with custom cred. provider
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(provider))
	if err != nil {
		log.Fatalf("could not load config. %s", err.Error())
	}
	if yndxConfig.Debug {
		logger.Println("[YANDEX S3] initialized configs and credential providers")
	}
	// creating s3 client
	s3client := s3.NewFromConfig(cfg)

	if yndxConfig.Debug {
		logger.Println("[YANDEX S3] initialized s3 Client")
	}

	// creating wrapper-client
	client := &YandexS3Client{
		client: s3client,
		owner:  yndxConfig.Owner,
		debug:  yndxConfig.Debug,
		bucket: yndxConfig.Bucket,
		logger: logger,
	}

	if yndxConfig.Debug {
		logger.Println("[YANDEX S3] initialized Yandex S3 Client")
	}
	return client
}
