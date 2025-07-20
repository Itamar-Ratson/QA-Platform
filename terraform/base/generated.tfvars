vpc_cidr = "10.0.0.0/16"
private_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
public_subnets = ["10.0.101.0/24", "10.0.102.0/24"]
availability_zones = ["eu-north-1a", "eu-north-1b"]
environment = "qa-test"
region = "eu-north-1"
common_tags = {
  "TestTimestamp" = "2025-07-20T15:23:37Z"
  "CreatedBy" = "qa-test-app"
  "AutoCleanup" = "true"
  "Environment" = "test"
  "TestCase" = "VPC Connectivity Test"
  "TestWorkspace" = "test-vpc-connectivity-test-1753014213"
}
