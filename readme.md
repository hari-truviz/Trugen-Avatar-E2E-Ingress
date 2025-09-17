# Infra Ingress

A Go-based WebSocket service that manages conversation requests through a queuing system, interfacing with RunPod infrastructure for processing.

## Overview

Infra Ingress is a real-time WebSocket server that:
- Accepts conversation requests via WebSocket connections
- Queues requests for processing
- Integrates with RunPod infrastructure
- Provides health monitoring endpoints

## Features

- **WebSocket Support**: Real-time bidirectional communication
- **Request Queuing**: Thread-safe FIFO queue with connection-based indexing
- **RunPod Integration**: Automatic request processing when infrastructure is ready
- **Connection Management**: Automatic cleanup of disconnected clients
- **Health Monitoring**: Built-in health check endpoint

## Architecture

```
Client (WebSocket) → Request Queue → RunPod Infrastructure
                           ↓
                    Periodic Processing
```

The service uses:
- `container/list` for efficient queue operations
- `sync.Mutex` for thread-safe access
- Connection-based indexing for request tracking
- Periodic ticker for infrastructure polling

## Installation

### Prerequisites

- Go 1.16 or higher
- Access to RunPod infrastructure

### Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd infra-ingress
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the application:
```bash
go build -o infra-ingress
```

## Usage

### Starting the Server

```bash
./infra-ingress
```

The server will start on port 8080 with the following endpoints:
- WebSocket: `ws://localhost:8080/ws`
- Health Check: `http://localhost:8080/health`

### WebSocket API

#### Connecting

Connect to the WebSocket endpoint:
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');
```

#### Sending Requests

Send conversation requests as JSON:
```json
{
  // ConversationRequest fields
  // (structure depends on CustomTypes.ConversationRequest)
}
```

#### Response

Upon successful queuing, you'll receive an acknowledgment:
```
"Conversation request added into the queue."
```

## Configuration

The service uses the following configuration:
- **Port**: 8080 (hardcoded)
- **Processing Interval**: 1 second
- **WebSocket Buffer Sizes**: 1024 bytes (read/write)

## API Reference

### WebSocket Endpoints

| Endpoint | Protocol | Description |
|----------|----------|-------------|
| `/ws` | WebSocket | Main conversation request endpoint |

### HTTP Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Service health check |

## Request Flow

1. Client connects to WebSocket endpoint
2. Client sends conversation request (JSON)
3. Server queues the request
4. Server sends acknowledgment to client
5. Background process checks RunPod availability
6. When ready, request is dequeued and processed
7. Connection cleanup on disconnect

## Queue Management

The service implements a thread-safe queue with:
- **FIFO Processing**: First-in, first-out request handling
- **Connection Indexing**: Fast lookup and cleanup by connection
- **Automatic Cleanup**: Removes all requests when client disconnects

## Monitoring

### Health Check

```bash
curl http://localhost:8080/health
```

Response: `Trugen Infra Ingress running.`

### Queue Status

The service logs queue length periodically:
```
Number of requests in the queue: X
```

## Dependencies

### External Packages
- `github.com/gorilla/websocket` - WebSocket implementation

### Internal Packages
- `infra-ingress/customtypes` - Custom type definitions
- `infra-ingress/utilities/infrastructure/runpod` - RunPod utilities

## Development

### Project Structure

```
infra-ingress/
├── main.go                           # Main application
├── customtypes/                      # Custom type definitions
├── utilities/infrastructure/runpod/  # RunPod integration
└── README.md                         # This file
```

### Key Components

- **RequestQueue**: Thread-safe queue with connection indexing
- **Request**: Wrapper containing WebSocket connection and conversation request
- **handleConnection**: WebSocket message handler
- **Background Processor**: Periodic RunPod polling and request processing

## Error Handling

The service handles:
- WebSocket connection errors
- JSON unmarshaling errors
- Queue operation errors
- RunPod infrastructure errors

Errors are logged to stdout for monitoring.
