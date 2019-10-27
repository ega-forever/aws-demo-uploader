package bucket

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	AWSS3 "github.com/aws/aws-sdk-go/service/s3"
	AWSS3Manager "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	s3         *AWSS3.S3
	bucket     string
	downloader *AWSS3Manager.Downloader
}

func New(bucket string, region string) *S3 {

	creds := credentials.NewEnvCredentials()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region:      aws.String(region),
			Credentials: creds,
		},
	}))

	s3 := AWSS3.New(sess)
	downloader := AWSS3Manager.NewDownloader(sess)

	return &S3{
		s3:         s3,
		bucket:     bucket,
		downloader: downloader,
	}
}

func (s3 *S3) GetFile(filename string) (*[]byte, error) {

	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := s3.downloader.Download(buf,
		&AWSS3.GetObjectInput{
			Bucket: aws.String(s3.bucket),
			Key:    aws.String(filename),
		})
	if err != nil {
		return nil, err
	}

	bytes := buf.Bytes()
	return &bytes, nil
}
