package spotFleet

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/service/elbv2"
    "github.com/dannybrody/aws-sdk-go-base/sessionHelper"
    "github.com/dannybrody/aws-sdk-go-base/errorHelper"
    "fmt"
    "sync"
)

const(
    InstanceHealthHealthy = "healthy"
)

type SpotFleet struct {
    SpotFleetRequestsOutput *ec2.DescribeSpotFleetRequestsOutput
    Ec2Svc                  *ec2.EC2
    Elbv2Svc                *elbv2.ELBV2          
}

// get all spot fleets
func describeSpotFleetRequests() *ec2.DescribeSpotFleetRequestsOutput{
    svc := ec2.New(sessionHelper.GetAwsSession())
    input := &ec2.DescribeSpotFleetRequestsInput{}

    result, err := svc.DescribeSpotFleetRequests(input)
    errorHelper.HandleAwsError(err)
    return result
}

// Create a new spot fleet object 
func New() *SpotFleet {
    return &SpotFleet{
        SpotFleetRequestsOutput: describeSpotFleetRequests(),
        Ec2Svc: ec2.New(sessionHelper.GetAwsSession()),
        Elbv2Svc: elbv2.New(sessionHelper.GetAwsSession()),
    }
}

// returns a slice of fleets that are in the desired state
func (fleet *SpotFleet) FilterFleetRequestsByState(state string) []ec2.SpotFleetRequestConfig {
    var results []ec2.SpotFleetRequestConfig

    for _, v := range fleet.SpotFleetRequestsOutput.SpotFleetRequestConfigs {
        if v.SpotFleetRequestState != nil && *v.SpotFleetRequestState == state {
            results = append(results, *v)
        }
    }
    return results
}

// get a slice of target group arns for a spot fleet
func (fleet *SpotFleet) GetFleetTargetGroupArns (f ec2.SpotFleetRequestConfig) []string{
    var fleetArns []string
    for _,arn := range f.SpotFleetRequestConfig.LoadBalancersConfig.TargetGroupsConfig.TargetGroups{
        fleetArns = append(fleetArns, *arn.Arn)
    }
    return fleetArns
}



// get all instances in a spot fleet
func getSpotFleetInstances (fleetId *string) *ec2.DescribeSpotFleetInstancesOutput{
    svc := ec2.New(sessionHelper.GetAwsSession())
    result, err := svc.DescribeSpotFleetInstances(&ec2.DescribeSpotFleetInstancesInput{
        SpotFleetRequestId: fleetId,
    })
    errorHelper.HandleAwsError(err)
    return result
}

// external method to get instances in a spot fleet
func (fleet *SpotFleet) GetSpotFleetInstances(fleetId *string) *ec2.DescribeSpotFleetInstancesOutput{
    return getSpotFleetInstances(fleetId)
}

// function that will return true if all instances in the fleet are healthy
func (fleet *SpotFleet) AreAllFleetInstancesHealthy(fleetId *string) bool{
    fleetInstances := getSpotFleetInstances(fleetId)
    for _, activeInstance := range fleetInstances.ActiveInstances{
        if activeInstance.InstanceHealth != nil && *activeInstance.InstanceHealth != InstanceHealthHealthy {
            return false
        }
    }
    return true
}

// cancel a spot fleet
func (fleet *SpotFleet) CancelSpotFleetRequest(spotFleetRequestId *string){
    fmt.Println("Cancelling Spot Fleet request: ", *spotFleetRequestId)
    // svc := ec2.New(sessionHelper.GetAwsSession())
    input := &ec2.CancelSpotFleetRequestsInput{
        SpotFleetRequestIds: []*string{
            spotFleetRequestId,
        },
        TerminateInstances: aws.Bool(true),
    }

    result, err := fleet.Ec2Svc.CancelSpotFleetRequests(input)
    errorHelper.HandleAwsError(err)
    fmt.Println(result)
}

// deregister an instance from a spot fleet
func (fleet *SpotFleet) DeregisterInstanceFromTargetGroup(wg *sync.WaitGroup, instanceId *string , targetGroupArn string ){
    defer wg.Done()
    // svc := elbv2.New(sessionHelper.GetAwsSession())
    fmt.Println("deregisteringTarget", *instanceId, targetGroupArn)

    input := &elbv2.DeregisterTargetsInput{
        TargetGroupArn: aws.String(targetGroupArn),
        Targets: []*elbv2.TargetDescription{
            {
                Id: instanceId,
            },
        },
    }

    _, err := fleet.Elbv2Svc.DeregisterTargets(input)
    if err != nil {
        fmt.Println(err.Error())
    }
}
