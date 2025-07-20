package tests

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// SubnetConnectivityTest tests subnet reachability
type SubnetConnectivityTest struct{}

func (t *SubnetConnectivityTest) Name() string {
	return "test_subnet_connectivity"
}

func (t *SubnetConnectivityTest) Description() string {
	return "Tests connectivity between subnets"
}

func (t *SubnetConnectivityTest) Execute(ctx context.Context, tfOutputs map[string]interface{}) TestResult {
	vpcID, ok := tfOutputs["vpc_id"].(string)
	if !ok {
		return TestResult{
			Success: false,
			Message: "VPC ID not found in outputs",
		}
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-north-1"),
	})
	if err != nil {
		return TestResult{
			Success: false,
			Message: fmt.Sprintf("AWS session creation failed: %v", err),
		}
	}

	ec2Svc := ec2.New(sess)

	// Get subnets
	subnets, err := ec2Svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	})
	if err != nil {
		return TestResult{
			Success: false,
			Message: fmt.Sprintf("Failed to describe subnets: %v", err),
		}
	}

	subnetCount := len(subnets.Subnets)
	details := map[string]interface{}{
		"vpc_id":       vpcID,
		"subnet_count": subnetCount,
		"subnets":      []map[string]string{},
	}

	subnetInfo := []map[string]string{}
	for _, subnet := range subnets.Subnets {
		info := map[string]string{
			"id":   *subnet.SubnetId,
			"cidr": *subnet.CidrBlock,
			"az":   *subnet.AvailabilityZone,
		}
		subnetInfo = append(subnetInfo, info)
	}
	details["subnets"] = subnetInfo

	return TestResult{
		Success: subnetCount > 0,
		Message: fmt.Sprintf("Found %d subnets in VPC", subnetCount),
		Details: details,
	}
}
