package utils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"mime/multipart"
	"net/url"
	"strings"
)

var (
	s3Client   *s3.Client
	bucketName string
)

// InitializeS3 initializes the S3 client and bucketName.
func InitializeS3(client *s3.Client, bucket string) {
	s3Client = client
	bucketName = bucket
}

// UploadFileToS3 uploads a file to the specified S3 bucket and returns the file URL.
func UploadFileToS3(file multipart.File, fileName string) (string, error) {
	sanitizedFileName := strings.ReplaceAll(fileName, " ", "_")
	sanitizedFileName = url.QueryEscape(sanitizedFileName)
	fmt.Println("sanitized name: ", sanitizedFileName)
	// Set the content type and other metadata as needed
	buffer := make([]byte, 1024)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String("uploads/" + sanitizedFileName),
		Body:        file,
		ContentType: aws.String("image/jpeg"), // Adjust content type if necessary
	}

	_, err = s3Client.PutObject(context.TODO(), uploadInput)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/uploads/%s", bucketName, sanitizedFileName)
	return fileURL, nil
}

// DeleteFileFromS3 deletes a file from S3 by its key
func DeleteFileFromS3(fileName string) error {
	sanitizedFileName := strings.ReplaceAll(fileName, " ", "_")
	_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("uploads/" + sanitizedFileName),
	})
	if err != nil {
		log.Printf("Failed to delete file from S3: %v", err)
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}
	return nil
}
