# AGENTS.md - Development Guide for AI Agents

This file provides development guidelines and commands for AI agents working on the SwipeLearn project.

## Project Overview

SwipeLearn is a spaced repetition web application built to demonstrate DevOps and Platform Engineering skills. The project prioritizes operability over feature richness.

**Tech Stack:**
- Backend: Go (REST API)
- Frontend: React (mobile-first)
- Database: PostgreSQL
- Infrastructure: Kubernetes (k3d/k3s), Terraform
- CI/CD: GitHub Actions
- Containerization: Docker (multi-stage builds)

## Repository Structure (Planned)

```
swipelearn/
├── app/              # Go backend
├── frontend/         # React frontend
├── infra/            # Terraform (local & prod)
├── k8s/              # Kubernetes manifests
├── .github/workflows # CI pipelines
├── docs/             # Architecture & decisions
└── AGENTS.md         # This file
```

## Development Commands

### Go Backend (app/)
```bash
# Run tests
go test ./...

# Run single test
go test -run TestSpecificFunction ./path/to/package

# Run tests with verbose output
go test -v ./...

# Build application
go build -o bin/server ./cmd/server

# Run development server
go run ./cmd/server

# Format code
go fmt ./...

# Lint code (requires golangci-lint)
golangci-lint run

# Vendor dependencies
go mod tidy
go mod vendor
```

### React Frontend (frontend/)
```bash
# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build

# Run tests
npm test

# Run single test
npm test -- --testNamePattern="SpecificTest"

# Run tests with coverage
npm test -- --coverage

# Lint code
npm run lint

# Format code
npm run format

# Type check
npm run typecheck
```

### Terraform (infra/)
```bash
# Initialize Terraform
terraform init

# Plan changes
terraform plan

# Apply changes
terraform apply

# Destroy infrastructure
terraform destroy

# Validate configuration
terraform validate

# Format configuration
terraform fmt

# Lint configuration (requires tflint)
tflint
```

### Docker
```bash
# Build backend image
docker build -t swipelearn-backend -f app/Dockerfile .

# Build frontend image
docker build -t swipelearn-frontend -f frontend/Dockerfile .

# Build all images
docker-compose build

# Run local development
docker-compose up -d
```

### Kubernetes
```bash
# Apply manifests
kubectl apply -f k8s/

# Get pod status
kubectl get pods -n swipelearn

# View logs
kubectl logs -f deployment/swipelearn-api -n swipelearn

# Port forward
kubectl port-forward service/swipelearn-api 8080:80 -n swipelearn
```

## Code Style Guidelines

### Go Backend
- Follow standard Go formatting (`go fmt`)
- Use `gofmt` for all code formatting
- Package names: short, lowercase, single words when possible
- Exported functions must have comments
- Error handling: always handle errors, use fmt.Errorf for wrapping
- Naming: CamelCase for exported, camelCase for unexported
- Interface names: usually -er suffix (e.g., `Reader`, `Writer`)
- Use standard library first, then minimal external dependencies
- Struct tags: JSON tags must match API contract
- Logging: use structured logging (logrus or similar)

### React Frontend
- Use TypeScript for all new code
- Functional components with hooks
- Component naming: PascalCase
- File naming: PascalCase for components, camelCase for utilities
- Prefer explicit types over `any`
- Use CSS modules or styled-components, avoid inline styles
- Mobile-first responsive design
- Accessibility: use semantic HTML, ARIA labels where needed
- State management: prefer React Context for simple state, consider Redux for complex state

### General Guidelines
- Write clear, commit messages following conventional commits
- Keep functions small and focused
- Prefer readability over cleverness
- Use environment variables for configuration
- Never commit secrets or API keys
- All user-facing text should be in English
- Follow semantic versioning for releases
- Write tests for critical business logic
- Documentation should be kept up-to-date

## Import Organization

### Go
```go
import (
    "context"
    "fmt"
    "log"
    
    "github.com/gin-gonic/gin"
    "github.com/lib/pq"
    
    "swipelearn/internal/config"
    "swipelearn/internal/models"
)
```
Order: standard library, third-party, internal packages

### TypeScript/React
```typescript
import React, { useState, useEffect } from 'react';
import { Button, Card } from '@/components/ui';
import { apiClient } from '@/lib/api';
import { Flashcard } from '@/types';
import './FlashcardList.css';
```
Order: React imports, third-party libraries, internal imports, styles

## Error Handling

### Go
- Always check for errors
- Use explicit error variables with `var ErrSomething = errors.New("...")`
- Wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Log errors at appropriate levels
- Return meaningful HTTP status codes from API handlers

### React/TypeScript
- Use error boundaries for component errors
- Validate user inputs
- Handle API errors gracefully with user-friendly messages
- Use try-catch for async operations
- Don't throw in render methods

## Testing Guidelines

### Go
- Table-driven tests for multiple scenarios
- Mock external dependencies
- Test both success and failure paths
- Use subtests for related test cases
- Aim for >80% coverage on business logic

### React
- Unit tests for components using React Testing Library
- Integration tests for user workflows
- Mock API calls in tests
- Test accessibility
- Use meaningful test descriptions

## Database Guidelines

- Use migrations for schema changes
- All database operations should be transactional where appropriate
- Use prepared statements to prevent SQL injection
- Implement proper connection pooling
- Add indexes for frequently queried columns
- Follow naming conventions: snake_case for tables and columns

## API Design

- RESTful principles
- Use appropriate HTTP methods
- Consistent error response format
- API versioning when needed
- Request/response validation
- Rate limiting considerations
- Health check endpoints

## Security Considerations

- Validate all inputs
- Use parameterized queries
- Implement proper authentication/authorization
- Use HTTPS in production
- Set security headers
- Keep dependencies updated
- Regular security audits

## Performance Guidelines

- Optimize database queries
- Use caching where appropriate
- Monitor application performance
- Lazy load non-critical resources
- Optimize bundle size for frontend
- Use CDN for static assets in production

## Deployment Guidelines

- Use multi-stage Docker builds
- Non-root containers
- Health checks and readiness probes
- Graceful shutdown handling
- Environment-specific configurations
- Blue-green or canary deployments when possible

## Monitoring and Observability

- Structured logging with correlation IDs
- Metrics for business-critical operations
- Alert on error rates and latency
- Database connection pool monitoring
- Resource usage tracking

## Local Development Setup

1. Create local k3d cluster: `k3d cluster create dev --servers 1 --agents 1 --port "8080:80@loadbalancer"`
2. Install dependencies: Go, Node.js, Docker, kubectl, Terraform
3. Set up environment variables in `.env.local`
4. Initialize Terraform: `cd infra && terraform init`
5. Start development services with Docker Compose

## CI/CD Pipeline

GitHub Actions should:
- Run linting and formatting checks
- Execute test suites
- Build Docker images
- Push to container registry
- Deploy to staging/production environments

## Cost Considerations

- Keep production costs under $15/month
- Use resource limits in Kubernetes
- Optimize Docker image sizes
- Monitor cloud resource usage
- Consider serverless options for non-critical services