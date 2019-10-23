package errorHelper

import (
    "github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/service/elbv2"
    "log"
)

func HandleAwsError(err error) {

    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case elbv2.ErrCodeTargetGroupNotFoundException:
                log.Fatal(elbv2.ErrCodeTargetGroupNotFoundException, aerr.Error())
            default:
                log.Fatal(aerr.Error())
            }
        } else {
            log.Fatal(err.Error())
        }
    }
}
