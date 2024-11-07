# Wex Technical Implementation Project

Welcome to the Wex Technical Implementation Project! This project implements a scalable and robust REST API for handling financial transactions and exchange rate calculations. It adheres to the requirements outlined by the hiring team. It is built with Go, following Hexagonal Architecture (Ports and Adapters), ensuring high maintainability and testability.

## Overview

This API is designed to provide precise and efficient processing of exchange rate data and transactions, with a strong focus on modularity, extensibility, and high performance. The project includes:

- **Hexagonal Architecture**: Implements a clean separation of concerns, making the application easily adaptable to changes and improvements.
- **Unit and Integration Tests**: Ensures reliability and robustness by covering both isolated and integrated functionality.
- **Containerization**: Provides a Dockerfile for containerizing the application, enabling easy deployment and scalability.
- **Exchange Rate Calculations**: Processes exchange rates with data retrieved from external sources (Treasury Reporting Rates of Exchange API).
- **Transaction Management**: Handles transactions and provides mechanisms for validating and storing them.
- **Error Handling**: Robust error management for handling edge cases and ensuring a smooth user experience.

## Key Technologies

- **Language**:
   - **Go (Golang)**: The primary language used for building the API.

- **Libraries**:
   - **Chi**: A lightweight and idiomatic routing library for Go.
   - **Zerolog**: A fast, structured logging library for Go, allowing precise control over log levels and outputs.
   - **Jsoniter**: A high-performance JSON library for Go, used for efficient JSON serialization and deserialization.

- **Database**:
   - **BoltDB (bbolt)**: An embedded, key-value database for efficient data storage, used for persisting transaction data.

- **Configuration**:
   - **godotenv**: Manages environment variables with a `.env` file for configuration.

- **Testing**:
   - **Testify**: A testing toolkit for Go, providing additional assertions and mocks for testing.

## Project Structure

The project follows a modular structure, organized according to Hexagonal Architecture principles:

```text
wex-technical-implementation-project
├── cmd
│   └── main.go                                         # Entry point of the application
├── internal
│   ├── adapters                                    # Adapters layer for integrating external clients and repositories
│   │   ├── client
│   │   │   ├── treasury_exchange_rate.go               # External client for treasury exchange rates
│   │   │   ├── treasury_exchange_rate_errors.go        # Error handling for the treasury client
│   │   │   ├── treasury_exchange_rate_mock.go          # Mock client for testing
│   │   │   └── treasury_exchange_rate_test.go          # Tests for the treasury client
│   │   ├── handler
│   │   │   ├── http.go                                 # HTTP handler for API endpoints
│   │   │   ├── http_errors.go                          # Error handling for HTTP responses
│   │   │   └── http_test.go                            # Tests for HTTP handlers
│   │   └── repository
│   │       ├── boltdb.go                               # BoltDB repository implementation
│   │       ├── boltdb_errors.go                        # Error handling for BoltDB
│   │       └── boltdb_test.go                          # Tests for BoltDB repository
│   ├── core                                        # Core application layer (business logic)
│   │   ├── domain                                    # Domain layer containing core entities and models
│   │   │   ├── exchange_rate.go                        # Exchange rate domain model
│   │   │   ├── exchange_rate_errors.go                 # Error handling for exchange rate model
│   │   │   ├── exchange_rate_test.go                   # Tests for exchange rate domain model
│   │   │   ├── transaction.go                          # Transaction domain model
│   │   │   ├── transaction_errors.go                   # Error handling for transaction model
│   │   │   └── transaction_test.go                     # Tests for transaction domain model
│   │   ├── ports                                     # Ports defining interfaces for the adapters
│   │   │   ├── exchange_rate.go                        # Interface for exchange rate service
│   │   │   └── transaction.go                          # Interface for transaction service
│   │   └── services                                  # Service implementations for business logic
│   │       ├── transaction.go                          # Transaction service implementation
│   │       └── transaction_test.go                     # Tests for transaction service
├── .dockerignore                                   # Docker ignore file
├── .env.example                                    # Example environment file
├── .gitignore                                      # Git ignore file
├── Dockerfile                                      # Dockerfile for containerizing the app
├── go.mod                                          # Go module dependencies
├── Makefile                                        # Makefile for automation commands
└── README.md                                       # Project README file
```

## Getting Started

### Prerequisites

- **Golang**: Ensure Go 1.23.2 or higher is installed. Download it from the [official Go website](https://golang.org/dl/).
- **Docker**: To build and run the application in a containerized environment. [Get Docker here](https://docs.docker.com/get-docker/).
- **Make**: For running automation commands. Install it with `sudo apt install make` on Ubuntu, `brew install make` on macOS or `sudo yum install make` on RHEL-like distributions.

### Installation

1. Clone the repository:

    ```sh
    git clone github.com/dainfoo/wex-technical-implementation-project
    cd wex-technical-implementation-project
    ```

2. Copy the example environment file and configure the environment variables:

    ```sh
    cp .env.example .env
    ```

3. Run the application:

    ```sh
    make run/dev
    ```
### API Call

1. Save a new transaction (run in port 8080):

   ```sh
   curl -X POST http://localhost:8080/transactions \
      -H "Content-Type: application/json" \
      -d '{
            "description": "Sample Transaction",
            "timestamp": "2023-11-06T15:04:05Z",
            "amount_in_usd": 28.745
          }'
   ```

2. Retrieve the transaction by ID:

    ```sh
    curl -X GET http://localhost:8080/transactions/ID-FROM-THE-PREVIOUS-CALL/YOUR-CURRENCY-NAME
    ```
