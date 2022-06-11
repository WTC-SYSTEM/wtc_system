package aws

import (
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/internal/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Aws struct {
	Config  config.AwsConfig
	S3      *s3.S3
	session *session.Session
}

type ClientAws interface {
}

func NewAwsConfig(cfg config.AwsConfig) *aws.Config {
	// create credentials
	creds := credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, "")
	// create Config
	c := aws.NewConfig().WithRegion(cfg.Region).WithCredentials(creds)
	return c
}

func NewSession(cfg config.AwsConfig) (*session.Session, error) {
	c := NewAwsConfig(cfg)
	// create session
	newSession, err := session.NewSession(c)
	if err != nil {
		return &session.Session{}, err
	}
	return newSession, nil
}

func NewS3(cfg config.AwsConfig) (Aws, error) {
	c := NewAwsConfig(cfg)
	s, err := NewSession(cfg)
	if err != nil {
		return Aws{}, err
	}
	return Aws{
		Config:  cfg,
		S3:      s3.New(s, c),
		session: s,
	}, nil
}
