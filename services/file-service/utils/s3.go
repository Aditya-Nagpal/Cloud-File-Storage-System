package utils

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"

	ConfigEnv "github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/file-service/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Uploader struct {
	Client     *s3.Client
	BucketName string
	Region     string
}

func NewS3Uploader() (*S3Uploader, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(ConfigEnv.AppConfig.AWSRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			ConfigEnv.AppConfig.AWSAccessKeyId,
			ConfigEnv.AppConfig.AWSSecretAccessKey,
			"",
		)),
	)
	if err != nil {
		panic("Failed to load AWS config: " + err.Error())
	}

	client := s3.NewFromConfig(cfg)

	return &S3Uploader{
		Client:     client,
		BucketName: ConfigEnv.AppConfig.BucketName,
		Region:     ConfigEnv.AppConfig.AWSRegion,
	}, nil
}

func (u *S3Uploader) UploadFile(file multipart.File, fileHeader *multipart.FileHeader, key string) error {
	defer file.Close()

	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	_, err = u.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(u.BucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
		ACL:         types.ObjectCannedACLPublicRead, // For public read access
	})
	return err
}

func (u *S3Uploader) UploadFolder(key string) error {
	_, err := u.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(u.BucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte{}), // empty content
	})
	return err
}

func (u *S3Uploader) GetS3URL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.BucketName, u.Region, key)
}
