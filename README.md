# Ledger

A financial ledger system built for managing accounts and financial transactions

## Features

- Account management
- Transaction processing
- Ledger operations
- Bank statement

## Architecture

### Services
- **Account Service**: Handles account management and related operations  
  - Manage accounts

- **Ledger Service**: Manages financial transactions, bank statement and ledger operations
  - Processes financial transactions
  - Maintains transaction history
  - Generates bank statements
  - Implements double-entry bookkeeping

### Databases
- Account Database
- Ledger Database

## Technology Stack

### Backend Services
- Golang
- Libs:
  - gRPC - For inter-service communication between account and ledger
  - Chi - HTTP router
  - Viper - Configuration management  
  - Testify - Testing framework

### Infrastructure
- **Containerization**: Docker
- **Orchestration**: Kubernetes
- **Load Balancing**: Nginx
- **Database**: PostgreSQL for account and ledger services
- **Nginx**: Act as a reverse proxy and load balancer

## How to run

### Running locally
  1. [Run the application locally](/k8s-manifests/README.md)
  2. [Postman http requests](/docs/assets/ledger%20http%20requests.postman_collection.json)
