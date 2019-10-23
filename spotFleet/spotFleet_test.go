
package spotFleet

import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/ec2"
    "testing"
    "github.com/stretchr/testify/assert"
)


// returns a slice of fleets that are in the desired state
func TestFleet_FilterFleetRequestsByState(t *testing.T){

    fleet := &SpotFleet{
        SpotFleetRequestsOutput: mockdescribeSpotFleetRequests(),
    }
    results := fleet.FilterFleetRequestsByState("cancelled")

    assert.Equal(t, len(results), 1)
}

func mockdescribeSpotFleetRequests() *ec2.DescribeSpotFleetRequestsOutput {

	instanceTypes := [...]string{"t2.small,", "t2.medium"}

	fmt.Println(len(instanceTypes))
	op := &ec2.DescribeSpotFleetRequestsOutput{
		SpotFleetRequestConfigs: []*ec2.SpotFleetRequestConfig{
			{
				ActivityStatus: aws.String("fulfilled"),
				SpotFleetRequestId: aws.String("sfr-2c378267-f745-4a10-9b59-055b9a3d34d0"),
      			SpotFleetRequestState: aws.String("cancelled"),
      			SpotFleetRequestConfig: &ec2.SpotFleetRequestConfigData{
      				AllocationStrategy: aws.String("diversified"),
      				ClientToken: aws.String("terraform-20191016193542889300000001"),
      				ExcessCapacityTerminationPolicy: aws.String("Default"),
      				FulfilledCapacity: aws.Float64(0),
      				IamFleetRole: aws.String("arn:aws:iam::482201535275:role/aws-ec2-spot-fleet-tagging-role"),
        			InstanceInterruptionBehavior: aws.String("terminate"),

        			LaunchSpecifications: []*ec2.SpotFleetLaunchSpecification{
        				{},
        			},

        			LoadBalancersConfig: &ec2.LoadBalancersConfig{
        				TargetGroupsConfig: &ec2.TargetGroupsConfig{
        					TargetGroups: []*ec2.TargetGroup{
        						{
        							Arn: aws.String("arn:aws:elasticloadbalancing:us-east-1:482201535275:targetgroup/api-test03-443-tg/9db4a68a25bcc843"),
        						},
        					},
        				},
        			},

			        OnDemandAllocationStrategy: aws.String("lowestPrice"),
			        OnDemandFulfilledCapacity: aws.Float64(0),
			        OnDemandTargetCapacity: aws.Int64(0),
			        ReplaceUnhealthyInstances: aws.Bool(true),
			        TargetCapacity: aws.Int64(1),
			        TerminateInstancesWithExpiration: aws.Bool(false),
			        Type: aws.String("maintain"),
			        // ValidUntil: "2022-08-17 06:59:26 +0000 UTC",
      			},
				// CreateTime: "2019-10-16 19:35:43.251 +0000 UTC",
			},
		},
	}

    return op
}
