package config

import (
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	AccessKeyID     string
	SecretAccessKey string
	MyRegion        string
	Session         *session.Session
}

func NewS3(accessKeyID, secretAccessKey, myRegion string) S3 {
	e := S3{accessKeyID, secretAccessKey, myRegion, nil}
	return e
}

func (s3 *S3) ConnectAws() {
	session, err := session.NewSession(
		&aws.Config{
			Region: aws.String(s3.MyRegion),
			Credentials: credentials.NewStaticCredentials(
				s3.AccessKeyID,
				s3.SecretAccessKey,
				"",
			),
		})
	s3.Session = session
	if err != nil {
		panic(err)
	}
}

func (s3 S3) UploadImage(bucketName string, fileName string, file multipart.File) (string, error) {
	uploader := s3manager.NewUploader(s3.Session)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		ACL:    aws.String("public-read"),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		return "", err
	}

	//filepath := "https://" + bucketName + "." + "s3-" + s3.MyRegion + ".amazonaws.com/" + fileName
	return result.Location, nil
}
