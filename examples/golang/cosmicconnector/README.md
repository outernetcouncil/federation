## Cosmic Connector Example Implementation

This example showcases a basic setup of the Cosmic Connector, a Federation server implementation allowing for custom handlers per Federation RPC. The example demonstrates configuration, server setup, and interaction with gRPC services.

## Sample Configuration

The Cosmic Connector is configured using a [ConnectorParams](../../pkg/go/config/config.proto) message defined in Protocol Buffers. Below is an example of a configuration in text protobuf format:

```textproto
port: 50052
observability_params {
  channelz_address: "0.0.0.0:50051"
  pprof_address: "0.0.0.0:6060"
}
```

### Configuration Breakdown

- **Transport Security**:
  - `insecure: {}`: Indicates that the connector will communicate over plain-text HTTP/2. For secure deployments, consider using TLS by configuring `system_cert_pool` instead.
  - For more details, refer to [config/config.proto](../../pkg/go/config/config.proto).

- **Provider Endpoint URI**:
  - `provider_endpoint_uri`: Specifies the RESTful provider endpoint that the Cosmic Connector bridges.

- **Authentication Strategy**:
  - `auth_strategy { none: {} }`: Sets the authentication strategy to none. For secure communication, refer to the [serverauth](../../pkg/go/auth/serverauth.go) package to implement JWT-based authentication.

- **Observability Parameters**:
  - `channelz_address`: Address for the gRPC Channelz introspection service.
  - `pprof_address`: Address for the pprof HTTP server for profiling.

For detailed configuration options, see [config/config.proto](config/config.proto).

## Running the Example

The example can be built and run using Bazel.

### Prerequisites

- **Bazelisk**: We recommend using Bazelisk (the Bazel wrapper) instead of installing Bazel directly:
  - Installation guide: https://github.com/bazelbuild/bazelisk#installation
  - Bazelisk automatically downloads and uses the correct Bazel version
  - Current project recommends Bazel 7.4.0 or later
- **Dependencies**: All dependencies are managed through Bazel's MODULE.bazel file.
  - github.com/rs/zerolog - Logging
  - google.golang.org/grpc - gRPC functionality
  - google.golang.org/protobuf - Protocol Buffers support

### Building the Example

Build the example using Bazel:

```bash
bazel build //examples/golang/cosmicconnector:cosmic_connector
```

### Running the Example

Execute the built binary with the necessary flags:

```bash
bazel run //examples/golang/cosmicconnector:cosmic_connector -- \
  -config path/to/config.textproto \
  -log-level info
```

#### Flags:

- `-config`: Path to the text protobuf configuration file (e.g., `config.textproto`). Required.
- `-log-level`: Sets the logging level (options: `disabled`, `warn`, `panic`, `info` (default), `fatal`, `error`, `debug`, `trace`).
- `-dry-run`: (Optional) Validate the configuration without starting the server. Exits with a non-zero return code if the config is invalid.

**Example: Dry Run Configuration Validation**

```bash
bazel run //examples/golang/cosmicconnector:cosmic_connector -- \
  -config path/to/config.textproto -dry-run
```

For more details, refer to the [main.go](./main.go) source file.

## Sample gRPC Calls

Once the Cosmic Connector is running, you can interact with its gRPC services using `grpcurl`. Below are sample commands for various operations.

### ListServiceOptions

Retrieve service options based on specific filters.

```bash
grpcurl -plaintext -d '{
  "service_attributes_filters": {
    "bidirectional_service_attributes_filter": {
      "bandwidth_bps_minimum": 1000000,
      "one_way_latency_maximum": "100ms"
    }
  }
}' localhost:50052 outernet.federation.v1alpha.Federation/ListServiceOptions
```

*Refer to [grpc_server.go](../../pkg/go/server/grpc_server.go) and [handler.go](./handler/handler.go) for implementation details.*

### ScheduleService

Schedule a new service with specified parameters.

```bash
grpcurl -plaintext -d '{
  "requestor_edge_to_requestor_edge": {
    "x_interconnection_points": [{
      "uuid": "point1",
      "coordinates": {
        "entry": [{
          "interval": {
            "start_time": "2024-01-01T00:00:00Z",
            "end_time": "2024-01-02T00:00:00Z"
          },
          "geodetic_wgs84": {
            "longitude_deg": 0,
            "latitude_deg": 0,
            "height_wgs84_m": 0
          }
        }]
      }
    }],
    "y_interconnection_points": [{
      "uuid": "point2",
      "coordinates": {
        "entry": [{
          "interval": {
            "start_time": "2024-01-01T00:00:00Z",
            "end_time": "2024-01-02T00:00:00Z"
          },
          "geodetic_wgs84": {
            "longitude_deg": 1,
            "latitude_deg": 1,
            "height_wgs84_m": 0
          }
        }]
      }
    }],
    "directionality": "DIRECTIONALITY_BIDIRECTIONAL"
  },
  "service_attributes_filters": {
    "bidirectional_service_attributes_filter": {
      "bandwidth_bps_minimum": 1000000,
      "one_way_latency_maximum": "1s"
    }
  },
  "priority": 1
}' localhost:50052 outernet.federation.v1alpha.Federation/ScheduleService
```

*See [handler.go](./handler/handler.go) for the implementation of `ScheduleService`.*

### MonitorServices

Monitor existing services by their IDs.

```bash
grpcurl -plaintext -d '{"add_service_ids": ["service-1", "service-2"]}' localhost:50052 outernet.federation.v1alpha.Federation/MonitorServices
```

*Implemented in [handler.go](./handler/handler.go).*

### CancelService

Cancel an active service.

```bash
grpcurl -plaintext -d '{"service_id": "service-1"}' localhost:50052 outernet.federation.v1alpha.Federation/CancelService
```

*Refer to [handler.go](./handler/handler.go) for the `CancelService` functionality.*

## Project Structure

Understanding the project structure helps in navigating and customizing the Cosmic Connector.

```
/cosmicconnector/
├── BUILD.bazel     # Example binary build configuration
├── README.md       # This file
├── config/         # Configuration package
│   ├── BUILD.bazel # Build config for config package
│   ├── config.go
│   ├── config.proto
│   └── doc.go
├── handler/        # Example handler implementation
│   ├── BUILD.bazel # Build config for handler package
│   ├── handler.go
│   └── handler_test.go
└── main.go         # Example entry point
```

### Example Components

The example demonstrates:
1. Configuration loading and validation
2. Server setup and initialization
3. Handler implementation for Federation services
4. Authentication setup (disabled for simplicity)
5. Observability configuration
6. Testing implementation with handler_test.go

For implementation details of the core packages used in this example, see the [pkg/go README](../../pkg/go/README.md).
