package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var s3b *S3Bucket

type S3Bucket struct {
	BucketName    string
	Region        string
	Client        *s3.Client
	PresignClient *s3.PresignClient
}

func NewS3Bucket(bucketname, region string) *S3Bucket {
	s3b := &S3Bucket{
		BucketName: bucketname,
		Region:     region,
	}

	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Config
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKeyID, secretAccessKey, "",
		)),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3b.Client = s3.NewFromConfig(cfg)
	s3b.PresignClient = s3.NewPresignClient(s3b.Client)

	return s3b
}

func (s *S3Bucket) GeneratePresignedURL(objectKey string, contentType string, expiry time.Duration) (string, error) {
	req, err := s.PresignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.BucketName),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType), // MIME type of the file
	}, s3.WithPresignExpires(expiry))

	if err != nil {
		return "", err
	}

	return req.URL, nil
}

// Exists checks if the file exists in S3
func (s *S3Bucket) Exists(key string) (bool, error) {
	_, err := s.Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func prepareUpload(c *gin.Context) {
	var request struct {
		FileExtension string `json:"file_extension" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format.  Please provide {\"file_extension\": \"<extension>\"}"})
		return
	}

	fileExtension := request.FileExtension
	if fileExtension != ".jpeg" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file extension is not supported for upload"})
	}

	// Generate a UUID for the filename
	filename := uuid.New().String()
	objectKey := fmt.Sprintf("uploads/%s%s", filename, fileExtension) // Include extension in object key

	// Generate pre-signed URL
	url, err := s3b.GeneratePresignedURL(objectKey, "image/jpeg", 5*time.Minute)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate pre-signed URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"presigned_url": url,
		"key":           objectKey,
	})
}

func uploadConfirm(c *gin.Context) {
	var request struct {
		Key string `json:"key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if the file exists in S3
	ok, err := s3b.Exists(request.Key)
	if err != nil || !ok {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "File not found in S3"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "status": "file uploaded successfully"})
}

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get values from environment variables
	bucketName := os.Getenv("S3_BUCKET_NAME")
	region := os.Getenv("AWS_REGION")

	// Create a new bucket client
	s3b = NewS3Bucket(bucketName, region)

	r := gin.Default()

	// Register endpoints

	// prepare-upload returns a pre-signed URL
	r.POST("/prepare-upload", prepareUpload)

	// upload-confirm checks if file was indeed uploaded by client on S3
	r.POST("/upload-confirm", uploadConfirm)

	port := "8080"

	log.Printf("Starting server on port %s...", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
