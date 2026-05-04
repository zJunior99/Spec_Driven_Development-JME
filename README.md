# Spec Driven Development in Go

Demonstration project for the article **"Spec Driven Development in Go: The Contract First, Then the Code"** published on Medium.

Implements a payment processing API using OpenAPI 3.0 as specification and Ginkgo as behavior spec framework.

## Requirements

- Go 1.23 or higher

## Run

```bash
go run ./cmd/...
```

## Test

```bash
go test ./... -v
```

## Project Structure

```
├── api/
│   └── openapi.yaml          # API specification (contract)
├── cmd/
│   └── main.go               # Entry point
├── internal/
│   ├── payment/
│   │   ├── payment.go        # Domain & service
│   │   ├── payment_suite_test.go  # Ginkgo suite
│   │   └── payment_spec_test.go   # Behavior specs
│   └── repository/
│       └── inmemory.go       # In-memory repository
└── .github/workflows/
    └── ci.yml                # CI/CD pipeline
```

## Article

[Spec Driven Development in Go: The Contract First, Then the Code](#)
