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

	fullDestPath := inp.Destination + inp.FileName

	_, err = y.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:              &y.bucket,
		Key:                 &fullDestPath,
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

func (y *YandexS3Client) DeleteFile(ctx context.Context, inp *DeleteFileInput) error {

	fullPath := inp.Destination + inp.FileName

	_, err := y.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket:              &y.bucket,
		Key:                 &fullPath,
		ExpectedBucketOwner: &y.owner,
	})
	if err != nil {
		return err
	}
	y.logger.Println(fmt.Sprintf("[YANDEX S3] delete file '%s' from '%s'. bucket - %s", inp.FileName, fullPath, y.bucket))
	return nil
}

func (y *YandexS3Client) GetFiles(ctx context.Context) (*Storage, error) {

	result, err := y.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:              &y.bucket,
		ExpectedBucketOwner: &y.owner,
	})

	if err != nil {
		return nil, err
	}

	data := make(chan *File, len(result.Contents))
	storage := NewStorage()
	mimetypes := []string{PNG, JPG, TTF}

	go func(data chan *File) {
		for mimeIdx, mime := range mimetypes {
			for _, r := range result.Contents {
				dest := *r.Key
				name, ok := GetFileNameByExt(dest, mime)
				if !ok {
					continue
				}
				file := &File{
					Name:        name,
					Extension:   mime,
					Size:        r.Size,
					Destination: dest,
				}
				data <- file
			}
			if mimeIdx == len(mimetypes)-1 {
				close(data)
			}
		}
	}(data)

	for f := range data {
		switch f.Extension {
		case TTF:
			storage.Fonts = append(storage.Fonts, f)
		case PNG:
			storage.Images = append(storage.Images, f)
		case JPG:
			storage.Images = append(storage.Images, f)
		}
	}
	y.logger.Println(fmt.Sprintf("[YANDEX S3] get all files from bucket - %s", y.bucket))
	return storage, nil

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
