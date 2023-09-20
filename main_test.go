package emptys3bucket_test

import (
	"context"
	"errors"
	"fmt"
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
const unversionedBucketName = "emptys3bucket-unversioned-integration-test"
const versionedBucketName = "emptys3bucket-versioned-integration-test"

var s3Client *s3.Client

func setup() ([]func() error, error) {
	var cleanupFunctions []func() error

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		return cleanupFunctions, errors.New("could not setup config")
	}
	s3Client = s3.NewFromConfig(cfg)

	bucketsToCreate := []string{unversionedBucketName, versionedBucketName}

	for _, bucketName := range bucketsToCreate {
		deleteBucketFunction, err := createBucket(s3Client, bucketName)
		if deleteBucketFunction != nil {
			cleanupFunctions = append(cleanupFunctions, func() error { return deleteBucketFunction() })
		}
		if err != nil {
			return cleanupFunctions, err
		}
	}

	_, err = s3Client.PutBucketVersioning(context.TODO(), &s3.PutBucketVersioningInput{
		Bucket: aws.String(unversionedBucketName),
		VersioningConfiguration: &types.VersioningConfiguration{
			Status: types.BucketVersioningStatusEnabled,
		},
	})
	if err != nil {
		return cleanupFunctions, errors.New("could not apply versioning on bucket")
	}

	return cleanupFunctions, nil
}

func cleanup(cleanupFunctions []func() error) []error {
	var errs []error

	for _, function := range cleanupFunctions {
		err := function()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func TestMain(m *testing.M) {
	var exitCode = 0

	cleanupFunctions, err := setup()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "An error occurred in setup(): "+err.Error())
		exitCode = 1
	} else {
		exitCode = m.Run()
	}

	errs := cleanup(cleanupFunctions)

	if errs != nil {
		for _, err = range errs {
			_, _ = fmt.Fprintf(os.Stderr, "An error occurred in cleanup(): "+err.Error())
		}
		if exitCode == 0 {
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}

func TestEmptyBucketWithOneItemAndNoVersioning(t *testing.T) {
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(unversionedBucketName),
		Key:    aws.String("index.html"),
		Body:   strings.NewReader("<h1>Hello World</h1>"),
	})
	require.Nil(t, err)

	emptys3bucket.EmptyBucket(s3Client, unversionedBucketName)

	assertBucketIsEmpty(t, s3Client, unversionedBucketName)
}

func TestEmptyBucketWithTwoItemsAndNoVersioning(t *testing.T) {
	files := []string{"index.html", "index2.html"}
	for _, file := range files {
		_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(unversionedBucketName),
			Key:    aws.String(file),
			Body:   strings.NewReader("<h1>Hello World</h1>"),
		})
		require.Nil(t, err)
	}

	emptys3bucket.EmptyBucket(s3Client, unversionedBucketName)

	assertBucketIsEmpty(t, s3Client, unversionedBucketName)
}

func TestEmptyBucketWithVersioningEnabled(t *testing.T) {
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(versionedBucketName),
		Key:    aws.String("index.html"),
		Body:   strings.NewReader("<h1>Version 1</h1>"),
	})
	require.Nil(t, err)

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(versionedBucketName),
		Key:    aws.String("index.html"),
		Body:   strings.NewReader("<h1>Version 2</h1>"),
	})
	require.Nil(t, err)

	emptys3bucket.EmptyBucket(s3Client, versionedBucketName)

	assertBucketIsEmpty(t, s3Client, versionedBucketName)
}
