package errorHelper

import (
    "github.com/aws/aws-sdk-go/aws/awserr"
    "log"
)

func HandleAwsError(err error) {

    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            default:
                log.Fatal(aerr.Error())
            }
        } else {
            log.Fatal(err.Error())
        }
    }
}