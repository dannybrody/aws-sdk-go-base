package session

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "os"
)

func getAwsSession() *session.Session {
    return session.Must(session.NewSession(&aws.Config{
        Region: aws.String(os.Getenv("AWS_DEFAULT_REGION")),
    }))
}
