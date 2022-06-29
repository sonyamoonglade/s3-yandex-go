package s3yandex

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
)

const (
	ACCESS_ID  = "AWS_ACCESS_KEY_ID"
	SECRET_KEY = "AWS_SECRET_ACCESS_KEY"
)

const (
	ImageJPG = "image/jpeg"
	ImagePNG = "image/png"
	FontTTF  = "font/ttf"
)

var (
	// ErrSecretOrAccessKeyNotFound cant resolve process's env variables
	ErrSecretOrAccessKeyNotFound = awserr.New("EnvSecretOrAccessNotFound", "AWS_SECRET_ACCESS_KEY or AWS_ACCESS_KEY_ID not found in environment", nil)
	// grant giving access to read file
	grant = ""
)

type PutFileInput struct {
	FilePath    string
	FileName    string
	Destination string
	ContentType string
}

type YandexS3Config struct {
	Owner  string
	Bucket string
	Debug  bool
}
