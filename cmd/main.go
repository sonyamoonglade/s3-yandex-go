package main

import (
	context2 "context"
	"github.com/joho/godotenv"
	"github.com/sonyamoonglade/s3-yandex-go/pkg/s3yandex"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Cant read env vars. %s \n", err.Error())
	}
	envProvider := s3yandex.NewEnvCredentialsProvider()

	config := s3yandex.YandexS3Config{
		Owner:  os.Getenv("BUCKET_OWNER_ID"),
		Debug:  true,
		Bucket: "zharpizza-bucket",
	}
	context := context2.Background()
	client := s3yandex.NewYandexS3Client(envProvider, config)

	err := client.PutFile(context, &s3yandex.PutFileInput{
		FilePath:    "upload/",
		FileName:    "upload_me.jpg",
		Destination: "static/test",
		ContentType: s3yandex.ImageJPG,
	})

	if err != nil {
		log.Fatalf("Could not put file into s3. %s \n", err.Error())
	}

}
