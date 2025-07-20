# QA Test App - TODO List

## Project Overview
Go-based TUI app for QA testing that reads YAML test case metadata and Terraform tfvars data, providing clickable cards for environment provisioning and test execution.

## Core Features TODO

### 1. Project Setup
- [x] Initialize Go module (`go mod init qa-test-app`)
- [x] Set up basic project structure
- [x] Create `.gitignore` for Go projects
- [x] Add `README.md` with usage instructions

### 2. YAML Configuration
- [x] Define YAML schema for test cases with metadata + tfvars
- [x] Implement YAML parser for test metadata and Terraform variables
- [x] Create sample test case YAML files with embedded tfvars
- [ ] Add validation for required fields (name, type, priority, severity, expected_result)
- [ ] Validate Terraform variable structure

### 3. Terraform Integration
- [x] Create base Terraform modules for common AWS resources
- [x] Implement tfvars file generation from YAML
- [ ] Add Terraform execution wrapper (init, plan, apply, destroy)
- [ ] Handle Terraform state management
- [ ] Add environment cleanup on test completion
- [ ] Implement resource tagging for test identification

### 4. TUI Implementation
- [ ] Choose TUI library (bubbletea/tview)
- [ ] Create main menu interface
- [ ] Implement card-based layout for test cases
- [ ] Add card click handlers for environment provisioning
- [ ] Show Terraform execution status/progress
- [ ] Implement navigation between cards
- [ ] Add environment teardown buttons

### 5. Test Functions
- [ ] Create test function interface
- [ ] Implement CIDR range validation function
- [ ] Implement DNS resolution test function
- [ ] Add ping/connectivity tests against provisioned resources
- [ ] Create test result reporting
- [ ] Add AWS resource validation functions (EKS, VPC, etc.)

### 5. TDD Implementation
- [ ] Set up testing framework
- [ ] Write unit tests for YAML parser
- [ ] Write unit tests for test functions
- [ ] Write integration tests for TUI components
- [ ] Add test coverage reporting

### 6. DevOps Pipeline
- [ ] Create GitHub Actions workflow
- [ ] Add linting (golangci-lint)
- [ ] Add automated testing
- [ ] Add security scanning (gosec)
- [ ] Add build and release automation

### 7. Containerization
- [ ] Create Dockerfile for the app
- [ ] Create docker-compose.yml for easy deployment
- [ ] Add volume mounts for test configs
- [ ] Optimize Docker image size

### 8. Documentation
- [ ] Add inline code documentation
- [ ] Create usage examples
- [ ] Document test function development
- [ ] Document Terraform module usage
- [ ] Add troubleshooting guide

## Sample YAML Structure
```yaml
metadata:
  name: "VPC Connectivity Test"
  type: "network"
  priority: "high"
  severity: "critical"
  expected_result: "All subnets should be reachable"
  description: "Test VPC connectivity across AZs"

terraform:
  tfvars:
    region: "eu-north-1"
    vpc_cidr: "10.0.0.0/16"
    private_subnets: ["10.0.1.0/24", "10.0.2.0/24"]
    public_subnets: ["10.0.101.0/24", "10.0.102.0/24"]
    availability_zones: ["eu-north-1a", "eu-north-1b"]
    environment: "qa-test"

test_functions:
  - "validate_cidr_ranges"
  - "test_subnet_connectivity"
  - "verify_route_tables"
```

## File Structure
```
qa-test-app/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── tui/
│   ├── tests/
│   ├── terraform/
│   └── yaml/
├── terraform/
│   ├── modules/
│   │   ├── vpc/
│   │   ├── eks/
│   │   └── common/
│   └── base/
├── test-cases/
│   └── sample.yaml
├── Dockerfile
├── docker-compose.yml
├── .github/workflows/
│   └── ci.yml
└── go.mod
```

## Development Phases

### Phase 1: Core Functionality
- Basic YAML parsing (metadata + tfvars)
- Simple TUI with cards
- Basic Terraform execution
- Basic test functions

### Phase 2: Enhanced Features
- Advanced TUI interactions
- More test function types
- Terraform state management
- Error handling and rollback

### Phase 3: DevOps & Deployment
- CI/CD pipeline
- Docker setup
- Terraform modules
- Documentation

## Testing Strategy
- Unit tests for all components
- Integration tests for TUI flow
- E2E tests with sample configs
- Performance tests for large test suites

## Acceptance Criteria
- [ ] App reads YAML test configurations with metadata and tfvars
- [ ] TUI displays clickable test cards
- [ ] Cards trigger Terraform environment provisioning
- [ ] Test functions execute against provisioned infrastructure
- [ ] Environment cleanup after test completion
- [ ] Docker Compose deployment works
- [ ] CI/CD pipeline passes all checks
- [ ] 80%+ test coverage maintained
- [ ] Terraform state management handled properly
