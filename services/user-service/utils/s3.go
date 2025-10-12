package utils

import (
	"context"
	"fmt"
	"mime/multipart"

	ConfigEnv "github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/user-service/config"

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

func (u *S3Uploader) UploadDisplayPicture(file multipart.File, header *multipart.FileHeader, key string) (string, error) {
	defer file.Close()

	_, err := u.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(u.BucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(header.Header.Get("Content-Type")),
		ACL:         types.ObjectCannedACLPublicRead, // For public read access
	})

	return u.GetS3URL(key), err
}

func (u *S3Uploader) DeleteDisplayPicture(key string) error {
	paginator := s3.NewListObjectsV2Paginator(u.Client, &s3.ListObjectsV2Input{
		Bucket: aws.String(u.BucketName),
		Prefix: aws.String(key),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return err
		}

		for _, obj := range page.Contents {
			_, err := u.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
				Bucket: aws.String(u.BucketName),
				Key:    obj.Key,
			})
			if err != nil {
				return fmt.Errorf("failed to delete object %s: %w", *obj.Key, err)
			}
		}
	}

	return nil
}

func (u *S3Uploader) GetS3URL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.BucketName, u.Region, key)
}
