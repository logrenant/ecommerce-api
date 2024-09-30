package services

import (
	"context"
	"io"
	"time"

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
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	// Bucket kontrolü ve oluşturma
	err = client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := client.BucketExists(context.Background(), bucket)
		if errBucketExists == nil && exists {
		} else {
			return nil, err
		}
	}

	return &MinIOService{
		client: client,
		bucket: bucket,
	}, nil

}

func (s *MinIOService) UploadFile(ctx context.Context, file io.Reader, fileName string) (string, error) {
	// Dosyayı yükle
	_, err := s.client.PutObject(ctx, s.bucket, fileName, file, -1, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucket, fileName, 7*24*time.Hour, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

// UpdateFile güncelleme işlemi için yeni bir fonksiyon ekliyoruz.
func (s *MinIOService) UpdateFile(ctx context.Context, oldFileName string, newFile io.Reader, newFileName string) (string, error) {
	// Eski dosyayı sil
	if err := s.DeleteFile(ctx, oldFileName); err != nil {
		return "", err
	}

	// Yeni dosyayı yükle
	return s.UploadFile(ctx, newFile, newFileName)
}

func (s *MinIOService) DeleteFile(ctx context.Context, filename string) error {
	err := s.client.RemoveObject(ctx, s.bucket, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
