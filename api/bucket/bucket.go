package bucket

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func GetMinioConnection() (*minio.Client, error) {
	endpoint := os.Getenv("BUCKET_ENDPOINT")
	accessKeyID := os.Getenv("BUCKET_USER")
	secretAccessKey := os.Getenv("BUCKET_PASSWORD")
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return minioClient, nil

}

func CreatePresignedUrl(bucket string, path string) (string, error) {
	minioClient, err := GetMinioConnection()
	if err != nil {
		return "", err
	}

	url, err := minioClient.PresignedGetObject(context.Background(), bucket, path, time.Hour, nil)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return url.String(), nil
}
