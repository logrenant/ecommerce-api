package services

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOService struct {
	client *minio.Client
	bucket string
}

func NewMinIOService(endpoint, accessKey, secretKey, bucket string) (*MinIOService, error) {
	creds := credentials.NewStaticV4(accessKey, secretKey, "")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  creds,
		Secure: false, // Eğer HTTPS kullanıyorsanız true yapın
	})
	if err != nil {
		return nil, err
	}

	// Bucket kontrolü ve oluşturma
	err = client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := client.BucketExists(context.Background(), bucket)
		if errBucketExists == nil && exists {
			// Bucket zaten mevcut
		} else {
			return nil, err
		}
	}

	return &MinIOService{
		client: client,
		bucket: bucket,
	}, nil

}

// UploadFile dosyayı MinIO'ya yükler
func (s *MinIOService) UploadFile(ctx context.Context, file io.Reader, fileName string) error {
	// MinIO'ya dosyayı yükle
	_, err := s.client.PutObject(ctx, s.bucket, fileName, file, -1, minio.PutObjectOptions{})
	return err
}

func (s *MinIOService) DeleteFile(ctx context.Context, filename string) error {
	err := s.client.RemoveObject(ctx, s.bucket, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
