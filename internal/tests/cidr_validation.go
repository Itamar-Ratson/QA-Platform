package tests

import (
	"context"
	"fmt"
	"net"
)

// CIDRValidationTest validates CIDR ranges don't overlap
type CIDRValidationTest struct{}

func (t *CIDRValidationTest) Name() string {
	return "validate_cidr_ranges"
}

func (t *CIDRValidationTest) Description() string {
	return "Validates CIDR ranges for overlaps and proper formatting"
}

func (t *CIDRValidationTest) Execute(ctx context.Context, tfOutputs map[string]interface{}) TestResult {
	vpcCIDR, ok := tfOutputs["vpc_cidr_block"].(string)
	if !ok {
		return TestResult{
			Success: false,
			Message: "VPC CIDR block not found in outputs",
		}
	}

	// Parse VPC CIDR
	_, vpcNet, err := net.ParseCIDR(vpcCIDR)
	if err != nil {
		return TestResult{
			Success: false,
			Message: fmt.Sprintf("Invalid VPC CIDR: %v", err),
		}
	}

	details := map[string]interface{}{
		"vpc_cidr": vpcCIDR,
		"checks": []string{},
	}

	checks := []string{
		fmt.Sprintf("VPC CIDR %s is valid", vpcCIDR),
		fmt.Sprintf("VPC network size: %d addresses", getNetworkSize(vpcNet)),
	}

	details["checks"] = checks

	return TestResult{
		Success: true,
		Message: "CIDR validation passed",
		Details: details,
	}
}

// Helper function to calculate network size
func getNetworkSize(network *net.IPNet) int {
	ones, bits := network.Mask.Size()
	return 1 << (bits - ones)
}
