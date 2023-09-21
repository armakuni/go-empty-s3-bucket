package emptys3bucket_test

import (
	"testing"

	emptys3bucket "github.com/armakuni/go-empty-s3-bucket"
	"github.com/stretchr/testify/require"
)

const awsRegion = "eu-west-1"
const testBucketName = "emptys3bucket-integration-test"

func TestEmptyBucketWithOneItemAndNoVersioning(t *testing.T) {
	simpleS3Client, err := NewTestS3Client(awsRegion)
	require.Nil(t, err)

	err = simpleS3Client.CreateBucket(testBucketName)
	require.Nil(t, err)
	defer simpleS3Client.DeleteBucket(testBucketName)

	err = simpleS3Client.PutObject(testBucketName, "index.html", "<h1>Hello World</h1>")
	require.Nil(t, err)

	emptys3bucket.EmptyBucket(simpleS3Client.S3Client, testBucketName)

	assertBucketIsEmpty(t, simpleS3Client.S3Client, testBucketName)
	require.True(t, false)
}

func TestEmptyBucketWithTwoItemsAndNoVersioning(t *testing.T) {
	simpleS3Client, err := NewTestS3Client(awsRegion)
	require.Nil(t, err)

	err = simpleS3Client.CreateBucket(testBucketName)
	require.Nil(t, err)
	defer simpleS3Client.DeleteBucket(testBucketName)

	err = simpleS3Client.PutObject(testBucketName, "index.html", "<h1>Hello World</h1>")
	require.Nil(t, err)
	err = simpleS3Client.PutObject(testBucketName, "index2.html", "<h1>Hello World</h1>")
	require.Nil(t, err)

	emptys3bucket.EmptyBucket(simpleS3Client.S3Client, testBucketName)

	assertBucketIsEmpty(t, simpleS3Client.S3Client, testBucketName)
}

func TestEmptyBucketWithVersioningEnabled(t *testing.T) {
	simpleS3Client, err := NewTestS3Client(awsRegion)
	require.Nil(t, err)

	err = simpleS3Client.CreateVersionedBucket(testBucketName)
	require.Nil(t, err)
	defer simpleS3Client.DeleteBucket(testBucketName)

	err = simpleS3Client.PutObject(testBucketName, "index.html", "<h1>Version 1</h1>")
	require.Nil(t, err)
	err = simpleS3Client.PutObject(testBucketName, "index.html", "<h1>Version 2</h1>")
	require.Nil(t, err)

	emptys3bucket.EmptyBucket(simpleS3Client.S3Client, testBucketName)

	assertBucketIsEmpty(t, simpleS3Client.S3Client, testBucketName)
}
