package s3

import (
	"context"
	"easyflow-backend/src/api"
	"easyflow-backend/src/common"
	"easyflow-backend/src/enum"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func connect(logger *common.Logger, cfg *common.Config) (*s3.Client, error) {
	config, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.BucketAccessKeyId, cfg.BucketSecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		logger.PrintfError(err.Error())
		return nil, err
	}

	client := s3.NewFromConfig(config, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.BucketURL)
	})

	return client, nil
}

func GetObjects(logger *common.Logger, cfg *common.Config, bucketName string) (*s3.ListObjectsV2Output, *api.ApiError) {
	client, err := connect(logger, cfg)
	if err != nil {
		logger.PrintfError("An error happened while connecting to the bucket %s", bucketName)
		return nil, &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	listedObjects, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})
	if err != nil {
		logger.PrintfError("An error happened while listing objects in bucket %s", bucketName)
		return nil, &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	return listedObjects, nil
}

func GetDownloadURL(logger *common.Logger, cfg *common.Config, bucketName string, objectKey string) (*string, *api.ApiError) {
	client, err := connect(logger, cfg)
	if err != nil {
		logger.PrintfError("An error happened while connecting to the bucket %s", bucketName)
		return nil, &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	presigner := s3.NewPresignClient(client)

	req, err := presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 10 * time.Second
	})
	if err != nil {
		logger.PrintfError("Could not get object %s in bucket %s", objectKey, bucketName)
		return nil, &api.ApiError{
			Code:    http.StatusInternalServerError,
			Error:   enum.ApiError,
			Details: err,
		}
	}

	return &req.URL, nil
}
