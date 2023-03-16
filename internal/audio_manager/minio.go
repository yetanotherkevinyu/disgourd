package audiomanager

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func InitMinioConnection() *minio.Client {
	endpoint := "localhost:9000"
	accessKey := "minioaccess"
	secretKey := "miniosecret"
	useSSL := false
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to initialize minio client: %v", err)
	}
	return minioClient
}

func CreateBucket(client *minio.Client, bucketName string, location string) {
	err := client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}
	log.Printf("Bucket %s created successfully", bucketName)
}

func UploadFile(client *minio.Client, bucketName string, objectName string, buf *bytes.Buffer) {
	n, err := client.PutObject(context.Background(), bucketName, objectName, buf, int64(buf.Len()), minio.PutObjectOptions{})
	if err != nil {
		log.Fatalf("Failed to upload audio file to minio: %v", err)
	}
	log.Printf("Uploaded %v bytes to %s/%s", n, bucketName, objectName)
}

func GetFile(client *minio.Client, bucketName, string, objectName string) *minio.Object {
	n, err := client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalf("Failed to get audio file from minio: %v", err)
	}
	return n
}
