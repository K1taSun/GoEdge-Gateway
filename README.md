# GoEdge Gateway

![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)

GoEdge Gateway is a lightweight, high-performance IoT edge gateway written in Go. It acts as a bridge between MQTT sensor networks and persistent storage, providing a robust gRPC API for data access.

## 🏗 Architecture

The system consists of modular microservices:

*   **MQTT Ingestor**: Subscribes to MQTT topics, processes incoming sensor data, and forwards it to the Storage Service.
*   **Storage Service**: A gRPC server that handles data persistence to PostgreSQL and provides an API for querying readings.
*   **PostgreSQL**: Relational database for reliable storage of sensor time-series data.
*   **Mosquitto**: MQTT Broker for handling sensor communications.

```mermaid
graph LR
    Sensors[IoT Sensors] -->|MQTT| Mosquitto
    Mosquitto -->|MQTT| Ingestor[MQTT Ingestor]
    Ingestor -->|gRPC| Storage[Storage Service]
    Storage -->|SQL| DB[(PostgreSQL)]
    Client[External Client] -->|gRPC| Storage
```

## 🚀 Features

*   **High Performance**: Built with Go for low latency and high concurrency.
*   **Modular Design**: Separated ingestion and storage layers.
*   **Protocol Buffers**: Efficient serialization using gRPC.
*   **Data Persistence**: Reliable storage using PostgreSQL.
*   **Scalable**: Designed to handle high throughput of sensor readings.
*   **Batch Processing**: Support for batch operations to optimize database writes.

## 📋 Prerequisites

*   **Go** 1.24 or higher
*   **Docker** & **Docker Compose** (for infrastructure)
*   **Make** (optional, for using Makefile commands)
*   **Protoc** (optional, for regenerating gRPC code)

## 🛠 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/k1tasun/GoEdge-Gateway.git
cd GoEdge-Gateway
```

### 2. Start Infrastructure

Launch the required services (PostgreSQL and MQTT Broker) using Docker Compose:

```bash
make run-infra
# Or directly: docker-compose up -d postgres mqtt-broker
```

### 3. Build Services

Compile the binaries:

```bash
make build
```

This will create executables in the `bin/` directory:
*   `bin/mqtt-ingestor`
*   `bin/storage-service`

### 4. Run Services

Open two terminals to run the services:

**Terminal 1: Storage Service**
```bash
./bin/storage-service
```

**Terminal 2: MQTT Ingestor**
```bash
./bin/mqtt-ingestor
```

## ⚙️ Configuration

Services are configured via environment variables.

### Common / Storage Service
| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Port for gRPC Storage Service | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:password@localhost:5432/goedge?sslmode=disable` |

### MQTT Ingestor
| Variable | Description | Default |
|----------|-------------|---------|
| `MQTT_BROKER` | URL of the MQTT Broker | `tcp://localhost:1883` |
| `MQTT_TOPIC` | Topic pattern to subscribe to | `sensors/#` |

## 📡 API Documentation

The Storage Service exposes a gRPC API defined in `api/proto/gateway.proto`.

### Methods

*   `StoreReading(StoreReadingRequest)`: Store a single sensor reading.
*   `StoreBatch(StoreBatchRequest)`: Store multiple readings in a single transaction.
*   `GetReadings(GetReadingsRequest)`: Retrieve historical readings for a specific device.

### Example Reading Object
```json
{
  "device_id": "sensor-01",
  "type": "temperature",
  "value": 23.5,
  "unit": "celsius",
  "timestamp": "2023-10-27T10:00:00Z"
}
```

## 🧪 Development

### Running Tests

Run unit tests for all packages:

```bash
make test
```

### Regenerating Protobuf Files

If you modify `api/proto/gateway.proto`, regenerate the Go code:

```bash
make gen-proto
```

### E2E Testing

A simple end-to-end test client is provided to verify the flow:

```bash
go run cmd/e2e-test/main.go
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

