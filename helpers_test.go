package emptys3bucket_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

type SimpleS3Client struct {
	S3Client *s3.Client
	region   string
}

func NewTestS3Client(awsRegion string) (*SimpleS3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		return nil, errors.New("could not setup config")
	}
	s3Client := s3.NewFromConfig(cfg)

	return &SimpleS3Client{
		S3Client: s3Client,
		region:   awsRegion,
	}, nil
}

func (client SimpleS3Client) CreateBucket(bucketName string) error {
	_, err := client.S3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(client.region),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %s", bucketName, err)
	}
	return nil
}

func (client SimpleS3Client) CreateVersionedBucket(bucketName string) error {
	err := client.CreateBucket(bucketName)
	if err != nil {
		return err
	}

	_, err = client.S3Client.PutBucketVersioning(context.TODO(), &s3.PutBucketVersioningInput{
		Bucket: aws.String(testBucketName),
		VersioningConfiguration: &types.VersioningConfiguration{
			Status: types.BucketVersioningStatusEnabled,
		},
	})

	if err != nil {
		return fmt.Errorf("could not apply versioning on bucket %s: %s", bucketName, err)
	}

	return err
}

func (client SimpleS3Client) PutObject(bucketName string, key string, content string) error {
	_, err := client.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   strings.NewReader(content),
	})

	return err
}

func (client SimpleS3Client) DeleteBucket(bucketName string) error {
	_, err := client.S3Client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return fmt.Errorf("failed to delete bucket %s: %s", bucketName, err)
	}

	return nil
}

func assertBucketIsEmpty(t *testing.T, s3Client *s3.Client, bucketName string) {
	objects, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	require.Nil(t, err)
	require.Empty(t, objects.Contents)

	versions, err := s3Client.ListObjectVersions(context.TODO(), &s3.ListObjectVersionsInput{
		Bucket: aws.String(bucketName),
	})
	require.Nil(t, err)
	require.Empty(t, versions.Versions)
	require.Empty(t, versions.DeleteMarkers)
}
