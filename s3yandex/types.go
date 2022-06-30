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
	JPG      = ".jpg"
	PNG      = ".png"
	TTF      = ".ttf"
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

type PutFileWithBytesInput struct {
	ContentType string
	FileName    string
	Destination string
	FileBytes   *[]byte
}

type DeleteFileInput struct {
	FileName    string
	Destination string
}

type YandexS3Config struct {
	Owner  string
	Bucket string
	Debug  bool
}

type File struct {
	Name        string
	Extension   string
	Size        int64
	Destination string
}

type Storage struct {
	Images []*File
	Fonts  []*File
}

func NewStorage() *Storage {

	return &Storage{Images: []*File{}, Fonts: []*File{}}
}
