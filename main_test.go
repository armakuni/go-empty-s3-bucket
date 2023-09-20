package emptys3bucket_test

import (
	"context"
	"strings"
	"testing"

	emptys3bucket "github.com/armakuni/go-empty-s3-bucket"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/require"
)

// tests needed:
// - empty bucket with no versioning
// - empty bucket with versioning
// - empty bucket with versioning and delete markers
// - empty bucket with versioning and delete markers and delete all versions
const awsRegion = "eu-west-1"

func TestEmptyBucket(t *testing.T) {
	const bucketName = "emptys3bucket-integration-test"
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	require.Nil(t, err, "")
	svc := s3.NewFromConfig(cfg)

	_, err = svc.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(awsRegion),
		},
	})

	require.Nil(t, err)

	_, err = svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("index.html"),
		Body:   strings.NewReader("<h1>Hello World</h1>"),
	})
	require.Nil(t, err)
	emptys3bucket.EmptyBucket(svc, bucketName)
	_, err = svc.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})

	require.Nil(t, err)
}
