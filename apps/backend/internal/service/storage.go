package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// StorageService handles file operations with S3-compatible storage (MinIO/R2)
type StorageService struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

// NewStorageService creates a new StorageService
func NewStorageService(endpoint, accessKey, secretKey, bucket, publicURL string) (*StorageService, error) {
	if endpoint == "" {
		slog.Warn("S3 endpoint not configured, storage service disabled")
		return &StorageService{bucket: bucket, publicURL: publicURL}, nil
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               endpoint,
			HostnameImmutable: true,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return nil, fmt.Errorf("load S3 config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &StorageService{
		client:    client,
		bucket:    bucket,
		publicURL: publicURL,
	}, nil
}

// Upload uploads a file to S3 storage
func (s *StorageService) Upload(ctx context.Context, folder, filename string, reader io.Reader, contentType string) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("storage service not configured")
	}

	key := fmt.Sprintf("%s/%s", folder, filename)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &key,
		Body:        reader,
		ContentType: &contentType,
	})
	if err != nil {
		return "", fmt.Errorf("upload file: %w", err)
	}

	publicURL := s.publicURL
	if publicURL == "" {
		publicURL = fmt.Sprintf("http://localhost:9000/%s", s.bucket)
	}

	return fmt.Sprintf("%s/%s", publicURL, key), nil
}

// Delete removes a file from S3 storage
func (s *StorageService) Delete(ctx context.Context, key string) error {
	if s.client == nil {
		return fmt.Errorf("storage service not configured")
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("delete file: %w", err)
	}
	return nil
}

// EnsureBucket creates the bucket if it doesn't exist
func (s *StorageService) EnsureBucket(ctx context.Context) error {
	if s.client == nil {
		return nil
	}

	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &s.bucket,
	})
	if err != nil {
		// Bucket doesn't exist, create it
		_, err = s.client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: &s.bucket,
		})
		if err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}
		slog.Info("created S3 bucket", "bucket", s.bucket)
	}
	return nil
}

const maxAvatarSize = 5 * 1024 * 1024 // 5 MB

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// ValidateImage validates image file type and returns content type
func ValidateImage(data []byte) (string, error) {
	if len(data) > maxAvatarSize {
		return "", ErrFileTooLarge
	}

	contentType := http.DetectContentType(data)
	if _, ok := allowedImageTypes[contentType]; !ok {
		return "", ErrInvalidFileType
	}

	return contentType, nil
}

// GenerateAvatarKey generates a unique key for user avatar
func GenerateAvatarKey(userID uuid.UUID, contentType string) string {
	ext := allowedImageTypes[contentType]
	if ext == "" {
		ext = ".jpg"
	}
	return fmt.Sprintf("avatars/%s/avatar%s", userID.String(), ext)
}

// ExtractKeyFromURL extracts the S3 key from a full URL
func ExtractKeyFromURL(url, publicURL string) string {
	if publicURL != "" {
		return strings.TrimPrefix(url, publicURL+"/")
	}
	// Try to extract path after bucket name
	parts := strings.SplitN(url, "/", 5)
	if len(parts) >= 5 {
		return parts[4]
	}
	return filepath.Base(url)
}
