package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type S3ConnConfigStruct struct {
	Region         string
	RawFilesBucket string
	AccessKey      string
	SecretKey      string
	EndpointUrl    string
}

func (connConfig *S3ConnConfigStruct) GetS3Session() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:           aws.String(connConfig.Region),
		Endpoint:         aws.String(connConfig.EndpointUrl),
		Credentials:      credentials.NewStaticCredentials(connConfig.AccessKey, connConfig.SecretKey, ""),
		S3ForcePathStyle: aws.Bool(true),
	})
}
