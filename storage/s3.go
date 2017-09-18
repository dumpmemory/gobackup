package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/huacnlee/gobackup/config"
	"github.com/huacnlee/gobackup/logger"
	"os"
	"path"
	"path/filepath"
)

// Amazon S3 storage
type S3 struct {
	bucket string
	path   string
}

func (ctx *S3) perform(model config.ModelConfig, archivePath string) error {
	logger.Info("=> storage | Amazon S3")
	s3Viper := model.StoreWith.Viper
	s3Viper.SetDefault("region", "us-east-1")

	cfg := aws.NewConfig()
	endpoint := s3Viper.GetString("endpoint")
	if len(endpoint) > 0 {
		cfg.Endpoint = aws.String(endpoint)
	}
	cfg.Credentials = credentials.NewStaticCredentials(
		s3Viper.GetString("access_key_id"),
		s3Viper.GetString("secret_access_key"),
		s3Viper.GetString("token"),
	)
	cfg.Region = aws.String(s3Viper.GetString("region"))
	cfg.MaxRetries = aws.Int(s3Viper.GetInt("max_retries"))

	ctx.bucket = s3Viper.GetString("bucket")
	ctx.path = s3Viper.GetString("path")

	sess := session.Must(session.NewSession(cfg))
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", archivePath, err)
	}

	key := path.Join(ctx.path, filepath.Base(archivePath))

	input := &s3manager.UploadInput{
		Bucket: aws.String(ctx.bucket),
		Key:    aws.String(key),
		Body:   f,
	}

	logger.Info("-> S3 Uploading...")
	result, err := uploader.Upload(input)
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	logger.Info("=>", result.Location)
	return nil
}