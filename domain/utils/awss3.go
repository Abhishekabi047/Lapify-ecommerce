package utils

import (
	"fmt"
	"mime/multipart"
	"project/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func CreateSession(cfg config.S3Bucket) *session.Session {
	fmt.Println("cfg", cfg)
	sess := session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String(cfg.Region),
			Credentials: credentials.NewStaticCredentials(
				cfg.AccessKeyId,
				cfg.AccessKeySecret,
				"",
			),
		},
	))
	return sess
}

func UploadImageToS3(file *multipart.FileHeader, sess *session.Session) (string, error) {
	image, err := file.Open()
	if err != nil {
		return "", err
	}
	// fmt.Println("**", sess)
	defer image.Close()
	uploader := s3manager.NewUploader(sess)
	upload, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("lapify/producct image/"),
		Key:    aws.String(file.Filename),
		Body:   image,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}
	return upload.Location, nil
}
