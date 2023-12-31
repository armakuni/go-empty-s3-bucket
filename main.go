package emptys3bucket

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func EmptyBucket(svc *s3.Client, bucketName string) {
	deleteAllObjects(svc, bucketName)

	err := removeDanglingDeleteMarkers(svc, bucketName)
	if err != nil {
		fmt.Println("Could not remove delete markers:" + err.Error())
	}
}

func deleteAllObjects(svc *s3.Client, bucketName string) {
	var identifiers []types.ObjectIdentifier

	objects, err := svc.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		fmt.Println("ListObjectsV2 failed: " + err.Error())
	} else {
		for _, object := range objects.Contents {
			identifiers = append(identifiers, types.ObjectIdentifier{Key: object.Key})
		}
	}

	versions, err := svc.ListObjectVersions(context.TODO(), &s3.ListObjectVersionsInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		fmt.Println("ListObjectVersions failed: " + err.Error())
	} else {
		for _, version := range versions.Versions {
			identifiers = append(identifiers, types.ObjectIdentifier{Key: version.Key, VersionId: version.VersionId})
		}
	}

	_, err = svc.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{
			Objects: identifiers,
		},
	})

	if err != nil {
		fmt.Println("DeleteObjects failed:" + err.Error())
	}
}

func removeDanglingDeleteMarkers(svc *s3.Client, bucketName string) error {
	deleteMarkers, err := svc.ListObjectVersions(context.TODO(), &s3.ListObjectVersionsInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return err
	}

	var dmIdentifiers []types.ObjectIdentifier
	for _, deleteMarker := range deleteMarkers.DeleteMarkers {
		dmIdentifiers = append(dmIdentifiers, types.ObjectIdentifier{
			Key:       deleteMarker.Key,
			VersionId: deleteMarker.VersionId,
		})
	}

	if dmIdentifiers == nil {
		return nil
	}

	_, err = svc.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{
			Objects: dmIdentifiers,
		},
	})
	return err
}
