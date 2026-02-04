# HCP Server

hcp-server is the backend service for the Hybrid Consensus Platform (HCP) benchmarking tool. It manages benchmarks, collects metrics, handles transaction generation, and monitors node status.

## Features

- **Benchmark Management**: Create, list, and manage benchmark runs.
- **Transaction Processing**: Track and analyze transactions.
- **Node Monitoring**: Monitor node status, health, and metrics.
- **Metric Collection**: Real-time metric collection for analysis.
- **gRPC API**: High-performance API for internal and external communication.

## Tech Stack

- **Language**: Go 1.25
- **Framework**: gRPC
- **Database**: PostgreSQL (GORM)
- **Cache**: Redis (Planned)
- **Build Tool**: Make

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL
- Protobuf Compiler

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/fffeng99999/hcp-server.git
   cd hcp-server
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Generate Protobuf files:
   ```bash
   make proto
   ```

4. Build the server:
   ```bash
   make build
   ```

### Running

1. Set up configuration in `configs/config.yaml` (or use defaults).

2. Start the server:
   ```bash
   ./bin/hcp-server
   ```

3. Run with migrations:
   ```bash
   ./bin/hcp-server migrate
   ```

## Development

### Running Tests

```bash
make test
```

### Project Structure

- `api/proto`: Protobuf definitions
- `cmd/server`: Main entry point
- `internal/config`: Configuration management
- `internal/database`: Database connection
- `internal/grpc/handlers`: gRPC request handlers
- `internal/models`: Data models
- `internal/repository`: Data access layer
- `internal/service`: Business logic
- `scripts`: Utility scripts

## License

MIT
