# VPC Module

Simple VPC module for QA testing environments.

## File Structure

- `vpc.tf` - VPC and Internet Gateway
- `subnets.tf` - Public and private subnets  
- `route_tables.tf` - Route tables and associations
- `variables.tf` - Input variables
- `outputs.tf` - Resource outputs

## Features

- VPC with configurable CIDR
- Public and private subnets across multiple AZs
- Internet Gateway for public subnets
- Route tables with proper associations
- No NAT Gateway (keeping costs minimal for testing)

## Usage

```hcl
module "vpc" {
  source = "../modules/vpc"

  vpc_cidr           = "10.0.0.0/16"
  private_subnets    = ["10.0.1.0/24", "10.0.2.0/24"]
  public_subnets     = ["10.0.101.0/24", "10.0.102.0/24"]
  availability_zones = ["eu-north-1a", "eu-north-1b"]
  environment        = "qa-test"
}
```

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| vpc_cidr | CIDR block for VPC | string | yes |
| private_subnets | List of private subnet CIDRs | list(string) | yes |
| public_subnets | List of public subnet CIDRs | list(string) | yes |
| availability_zones | List of AZs | list(string) | yes |
| environment | Environment name | string | yes |

## Outputs

| Name | Description |
|------|-------------|
| vpc_id | VPC ID |
| vpc_cidr_block | VPC CIDR block |
| private_subnet_ids | Private subnet IDs |
| public_subnet_ids | Public subnet IDs |
| internet_gateway_id | Internet Gateway ID |

## Testing

Deploy using the base configuration:

```bash
cd terraform/base
terraform init
terraform plan -var-file="sample.tfvars"
terraform apply -var-file="sample.tfvars"
```
