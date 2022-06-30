package s3yandex

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"os"
)

type EnvProvider struct {
	retrieved bool
}

func NewEnvCredentialsProvider() *EnvProvider {
	return &EnvProvider{}
}
func (p *EnvProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {

	accessId := os.Getenv(ACCESS_ID)
	secretKey := os.Getenv(SECRET_KEY)
	if accessId == "" || secretKey == "" {
		return aws.Credentials{}, ErrSecretOrAccessKeyNotFound
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
