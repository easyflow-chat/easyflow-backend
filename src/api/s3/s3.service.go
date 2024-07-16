package s3

import (
	"bytes"
	"context"
	"easyflow-backend/src/common"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BucketBasics struct {
	S3Client *s3.Client
}

var bucketName = "easyflow"

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

func GetObjects(logger *common.Logger, cfg *common.Config) {
	client, err := connect(logger, cfg)
	if err != nil {
		return
	}

	listedBuckets, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})
	if err != nil {
		logger.PrintfError(err.Error())
		return
	}

	logger.PrintfInfo("Objects in bucket %s:", bucketName)
	for _, object := range listedBuckets.Contents {
		logger.PrintfInfo(*object.Key)
	}
}

func UploadFile(logger *common.Logger, cfg *common.Config, bucketName string, objectKey string, file *bytes.Buffer, fileName string) error {
	client, err := connect(logger, cfg)
	if err != nil {
		return err
	}
	openedFile := io.Reader(file)
	if err != nil {
		logger.PrintfError(err.Error())
	} else {
		_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			Body:   openedFile,
		})
		if err != nil {
			logger.PrintfError("Couldn't upload file %v to %v:%v. Here's why: %v\n", fileName, bucketName, objectKey, err)
			return err
		}
	}
	return err
}
