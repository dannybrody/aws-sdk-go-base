package sessionHelper

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "os"
)
const (
	awsDefaultRegion = "us-east-1"
)

func GetAwsSession() *session.Session {
	region := getEnv("AWS_DEFAULT_REGION", awsDefaultRegion)
    return session.Must(session.NewSession(&aws.Config{
        Region: aws.String(region),
    }))
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
