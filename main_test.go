package emptys3bucket_test

import (
	"context"
	"os"
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
const bucketName = "emptys3bucket-integration-test"

var svc *s3.Client

func setup() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		panic("could not setup config")
	}
	svc = s3.NewFromConfig(cfg)

	_, err = svc.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(awsRegion),
		},
	})

	if err != nil {
		panic("could not create bucket")
	}
}

func cleanup() {
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	cleanup()
	os.Exit(code)
}

func TestEmptyBucketWithOneItemAndNoVersioning(t *testing.T) {
	_, err := svc.PutObject(context.TODO(), &s3.PutObjectInput{
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

func TestEmptyBucketWithTwoItemsAndNoVersioning(t *testing.T) {

}
