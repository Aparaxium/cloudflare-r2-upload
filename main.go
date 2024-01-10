package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var client *s3.Client

const workingDir = "./"

func main() {

	log.Println("Starting...")

	ACCOUNT_ID := os.Args[1]
	ACCESS_KEY_ID := os.Args[2]
	SECRET_ACCESS_KEY := os.Args[3]
	BUCKET := os.Args[4]
	SOURCE_FILE := os.Args[5]
	DESTINATION_FILE := os.Args[6]

	// print directory files
	files, err := os.ReadDir(workingDir)
	if err != nil {
		log.Fatal("Error reading directory:", err)
	}
	for _, file := range files {
		log.Println(file.Name())
	}

	log.Println("Creating client...")
	// Setup endpoint
	r2Resolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", ACCOUNT_ID),
			}, nil
		})

	// Create configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(ACCESS_KEY_ID, SECRET_ACCESS_KEY, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal("Error creating config:", err)
	}

	// Create client
	client = s3.NewFromConfig(cfg)

	log.Println("Client created...")
	log.Println("Uploading file...")
	log.Println("Source file: " + os.Getenv("INPUT_SOURCE_FILE"))

	// Open the file to be uploaded
	file, err := os.Open(filepath.Join(workingDir, SOURCE_FILE))
	if err != nil {
		log.Fatal("Error opening file:", err)
	}

	// Upload object
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(BUCKET),
		Key:    aws.String(DESTINATION_FILE),
		Body:   file,
	})
	if err != nil {
		log.Fatal("Error putting object:", err)
	}
	log.Println("File uploaded...")

	log.Println("Finished...")
}
