package emptys3bucket_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func createBucket(s3Client *s3.Client, bucketName string) (func() error, error) {
	_, err := s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(awsRegion),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create bucket %s: %s", bucketName, err)
	}

	return func() error { return deleteBucket(s3Client, bucketName) }, nil
}

func deleteBucket(s3Client *s3.Client, bucketName string) error {
	_, err := s3Client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
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
