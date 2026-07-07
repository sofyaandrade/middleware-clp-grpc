# Industrial PLC Middleware

> **Work in progress:** this project is under active development as my undergraduate thesis (TCC). It is not production-ready yet.

A concurrent backend written in Go for integrating applications with industrial Programmable Logic Controllers (PLCs) over Modbus TCP.

The middleware exposes a REST API for managing PLCs and their tags while running background workers that continuously synchronize real-time industrial data. The project explores backend architecture, concurrency, industrial communication protocols, authentication, and role-based access control.

[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Status](https://img.shields.io/badge/status-in%20development-orange)](#project-status)
[![Protocol](https://img.shields.io/badge/protocol-Modbus%20TCP-00599C)](https://modbus.org/)

## Technical Highlights

- **Concurrent PLC workers:** each active PLC is managed by its own goroutine, allowing multiple devices to be monitored independently.
- **Channel-based command processing:** buffered Go channels serialize read and write commands for each PLC connection.
- **Graceful shutdown:** `context.Context`, cancellation functions, OS signal handling, and `sync.WaitGroup` coordinate the application lifecycle.
- **Thread-safe real-time cache:** `sync.Map` and mutex-protected maps store tag values and PLC connection statuses safely across goroutines.
- **Automatic reconnection:** workers detect connection failures and retry Modbus TCP communication without restarting the API.
- **Efficient batch operations:** contiguous tags are grouped into blocks to reduce the number of Modbus requests.
- **Hot reload of PLC workers:** changes to PLC or tag configuration trigger a non-blocking synchronization signal.
- **Layered backend design:** HTTP routes, controllers, use cases, repositories, domain models, and infrastructure concerns are separated.
- **API security:** JWT access and refresh tokens, bcrypt password hashing, and Casbin-based RBAC.

## Concurrency Model

The Modbus service maintains isolated runtime resources for every configured PLC:

1. A dedicated goroutine owns the PLC communication loop.
2. A buffered channel queues commands for that device.
3. A ticker periodically polls configured tags.
4. A cancellation function stops or reloads the worker when its configuration changes.
5. A shared `WaitGroup` ensures all workers finish during shutdown.

This design keeps device communication independent, prevents concurrent access to the same Modbus connection, and allows the service to add, reload, or remove PLC workers at runtime.

## Modbus Capabilities

- Modbus TCP master communication
- Automatic connection retry and connection-status tracking
- Read and write operations for:
  - Coils
  - Discrete inputs
  - Holding registers
  - Input registers
- Supported tag representations include `REAL`, `INT`, `BOOL`, and `DWORD`
- Configurable byte and word order conversion
- Block reads of up to 120 registers or 2,000 coils per request
- Previous-value fallback when an individual block read fails

## Tech Stack

| Area | Technology |
| --- | --- |
| Language | Go |
| HTTP framework | Gin |
| Industrial protocol | Modbus TCP |
| ORM | GORM |
| Database | SQLite |
| Authentication | JWT + bcrypt |
| Authorization | Casbin RBAC |
| Configuration | Viper |
| API documentation | Swagger tooling |
| Concurrency | Goroutines, channels, contexts, mutexes, `sync.Map`, `WaitGroup` |

## Project Structure

```text
.
|-- cmd/app/                         # Application startup and server configuration
|-- internal/
|   |-- api/
|   |   |-- controllers/             # HTTP request handling
|   |   |-- middlewares/             # JWT, RBAC, and CORS
|   |   `-- routes/                  # REST endpoint definitions
|   |-- domain/
|   |   |-- interfaces/              # Application contracts
|   |   |-- models/                  # Domain and persistence models
|   |   `-- security/                # Tokens and password hashing
|   |-- infrastructure/
|   |   |-- clp/                     # PLC worker lifecycle manager
|   |   |-- database/                # Database setup and migrations
|   |   |-- jobs/                    # Thread-safe real-time state
|   |   `-- modbusMaster/            # Modbus connection, polling, reads, and writes
|   |-- repository/                  # GORM data access
|   `-- usecase/                     # Application business rules
|-- go.mod
`-- main.go
```

## Running Locally

### Requirements

- Go 1.25 or newer
- A Modbus TCP device or simulator for PLC communication

### Setup

```bash
git clone https://github.com/sofyaandrade/middleware-clp-grpc.git
cd middleware-clp-grpc
go mod download
go run .
```

The API starts at `http://localhost:1710`. On first run, the application creates the SQLite database, runs migrations, and seeds the initial reference data.

PLC hardware is not required to start the API, but a reachable Modbus TCP server is required to test real-time reads and writes.

## Main API Resources

| Resource | Purpose |
| --- | --- |
| `/login` and `/refresh` | Authentication and token renewal |
| `/users` | User management |
| `/clps` | PLC configuration and connection status |
| `/tags` | Tag configuration and real-time values |
| `/swaps` | Available byte and word order options |
| `/type-clps` | PLC types |
| `/type-tags` | Supported tag data types |
| `/type-operations` | Supported Modbus operations |

Protected routes require a Bearer access token and are evaluated against Casbin authorization policies.

## Project Status

This repository is an **active undergraduate thesis project** and remains under development. The current implementation demonstrates the core architecture, API, authentication flow, database persistence, concurrent PLC lifecycle management, and Modbus TCP communication.

Planned improvements include:

- Automated unit and integration tests, including race detection
- Structured logging, metrics, and improved error reporting
- Environment validation and production configuration
- Complete OpenAPI documentation
- Docker-based development environment
- CI pipeline for formatting, tests, and static analysis
- Additional resilience and security hardening

## Academic Context

This project investigates how a Go backend can act as an integration layer between industrial automation equipment and higher-level applications. Its main focus is reliable concurrent communication with multiple PLCs while exposing device configuration and real-time process data through a secured HTTP API.
