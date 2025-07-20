package tests

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// RouteTableTest verifies route table configuration
type RouteTableTest struct{}

func (t *RouteTableTest) Name() string {
	return "verify_route_tables"
}

func (t *RouteTableTest) Description() string {
	return "Verifies route table configuration and associations"
}

func (t *RouteTableTest) Execute(ctx context.Context, tfOutputs map[string]interface{}) TestResult {
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

	// Get route tables
	routeTables, err := ec2Svc.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
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
			Message: fmt.Sprintf("Failed to describe route tables: %v", err),
		}
	}

	details := map[string]interface{}{
		"vpc_id":           vpcID,
		"route_table_count": len(routeTables.RouteTables),
		"route_tables":     []map[string]interface{}{},
	}

	routeTableInfo := []map[string]interface{}{}
	hasInternetRoute := false

	for _, rt := range routeTables.RouteTables {
		rtInfo := map[string]interface{}{
			"id":           *rt.RouteTableId,
			"associations": len(rt.Associations),
			"routes":       []string{},
		}

		routes := []string{}
		for _, route := range rt.Routes {
			if route.DestinationCidrBlock != nil {
				routeStr := *route.DestinationCidrBlock
				if route.GatewayId != nil && strings.Contains(*route.GatewayId, "igw-") {
					routeStr += " -> IGW"
					hasInternetRoute = true
				}
				routes = append(routes, routeStr)
			}
		}
		rtInfo["routes"] = routes
		routeTableInfo = append(routeTableInfo, rtInfo)
	}

	details["route_tables"] = routeTableInfo
	details["has_internet_route"] = hasInternetRoute

	return TestResult{
		Success: len(routeTables.RouteTables) > 0,
		Message: fmt.Sprintf("Found %d route tables, internet route: %v", 
			len(routeTables.RouteTables), hasInternetRoute),
		Details: details,
	}
}
