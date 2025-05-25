package bucket

import (
	"context"
	"io"
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
		return "", err
	}

	return url.String(), nil
}

func UploadFile(bucket string, path string, file io.Reader, size int64, contentType string) error {
	minioClient, err := GetMinioConnection()
	if err != nil {
		return err
	}

	info, err := minioClient.PutObject(context.Background(), bucket, path, file, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}
	log.Println(info)

	return nil
}
