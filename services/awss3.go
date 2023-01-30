package services

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

type AwsS3Service struct {
	Client   *s3.Client
	Uploader *manager.Uploader
}

func NewAwsS3Service() *AwsS3Service {
	serv := &AwsS3Service{}
	return serv
}

func (s *AwsS3Service) Init() error {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-east-1"))
	if err != nil {
		logrus.Fatalln("unable to load SDK config,", err)
	}

	client := s3.NewFromConfig(cfg)
	s.Client = client

	uploader := manager.NewUploader(client)
	s.Uploader = uploader

	return nil
}

func (s *AwsS3Service) Start() error {
	return nil
}
