package photo

import (
	b "bytes"
	"context"
	"fmt"
	"github.com/WTC-SYSTEM/logging"
	"github.com/WTC-SYSTEM/wtc_system/photo_service/pkg/client/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

type Storage interface {
	Create(ctx context.Context, bytes []byte, folder string) (string, error)
}

type storage struct {
	logger logging.Logger
	awsCfg aws.Aws
}

// NewStorage
func NewStorage(l logging.Logger, awsCfg aws.Aws) Storage {
	return &storage{
		logger: l,
		awsCfg: awsCfg,
	}
}

func (s storage) Create(ctx context.Context, bytes []byte, folder string) (string, error) {
	uniq, _ := uuid.NewUUID()
	fName := folder + "/" + uniq.String() + "-" + strconv.FormatInt(time.Now().Unix(), 10) + ".jpg"
	ct := http.DetectContentType(bytes)
	params := &s3.PutObjectInput{
		Bucket:      &s.awsCfg.Config.Bucket,
		Key:         &fName,
		Body:        b.NewReader(bytes),
		ContentType: &ct,
	}
	_, err := s.awsCfg.S3.PutObject(params)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.awsCfg.Config.Bucket, s.awsCfg.Config.Region, fName), nil
}
