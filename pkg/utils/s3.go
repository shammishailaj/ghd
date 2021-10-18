package utils

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// UploadToS3 - Function to upload a file to AWS' Simple Storage Service (S3)
// bucket - string containing the bucket name
// key - string containing the key under which the data is to be stored
// source - string containing the path of local file that is to be uploaded
// Cf. https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html
func (u *Utils) UploadToS3(bucket string, key string, source string) (*s3manager.UploadOutput, error) {
	sess, err := u.GetAWSSession()

	file, err := os.Open(source)
	if err != nil {
		u.ExitErrorf("Unable to open file %q, %v", err)
	}

	defer file.Close()

	// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
	// for more information on configuring part size, and concurrency.
	//
	// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
	uploader := s3manager.NewUploader(sess)

	uploadOutput, uploadErr := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if uploadErr != nil {
		// Print the error and exit.
		u.ExitErrorf("Unable to upload %q to %q, %v", key, bucket, err)
		log.Printf("Failed to upload %s to s3://%s/%s", source, bucket, key)
	}

	log.Printf("Successfully uploaded %q to %q\n", key, bucket)
	//log.Printf("Uploaded %s.tbz2 to %s with Upload ID %s and VersionID %s", source, uploadOutput.Location, uploadOutput.UploadID, *(uploadOutput.VersionID))
	log.Printf("Uploaded %s.tbz2 to %s with Upload ID %s and VersionID [REMOVED_FROM_CODE]", source, uploadOutput.Location, uploadOutput.UploadID)
	return uploadOutput, err
}
