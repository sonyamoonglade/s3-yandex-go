package main

import (
	"fmt"
)

func main() {
	//if err := godotenv.Load(); err != nil {
	//	log.Fatalf("Cant read env vars. %s \n", err.Error())
	//}
	//envProvider := s3yandex.NewEnvCredentialsProvider()
	//
	//config := s3yandex.YandexS3Config{
	//	Owner:  os.Getenv("BUCKET_OWNER_ID"),
	//	Debug:  true,
	//	Bucket: "zharpizza-bucket",
	//}
	//context := context2.Background()
	//client := s3yandex.NewYandexS3Client(envProvider, config)
	//files, err := client.GetFiles(context)
	//if err != nil {
	//	log.Fatalf("cant get files. %s", err.Error())
	//}
	//for _, image := range files.Images {
	//	fmt.Printf("'%s' from '%s'. Size - %d byte \n", image.Name, image.Destination, image.Size)
	//}
	//for _, font := range files.Fonts {
	//	fmt.Printf("'%s' from '%s'. Size - %d byte \n", font.Name, font.Destination, font.Size)
	//
	//}
	//
	//err := client.PutFile(context, &s3yandex.PutFileInput{
	//	FilePath:    "upload/",
	//	FileName:    "upload_me.jpg",
	//	Destination: "static/test/",
	//	ContentType: s3yandex.ImageJPG,
	//})
	//
	//if err != nil {
	//	log.Fatalf("Could not put file into s3. %s \n", err.Error())
	//}

	//err := client.DeleteFile(context, &s3yandex.DeleteFileInput{
	//	FileName:    "audi-rs6-avant-station-wagon-luxury-cars-3840x2560-5719.jpg",
	//	Destination: "static/check/",
	//})
	//if err != nil {
	//	log.Fatalf("Could not delete file from s3. %s \n", err.Error())
	//}

	fmt.Println(10 << 20) // (10 * 2 ^20) / (1024^2) megabyte

}
