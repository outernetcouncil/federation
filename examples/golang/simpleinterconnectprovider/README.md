## Interconnect Provider Example Implementation

This example showcases a basic setup of the Simple Interconnect Provider, a Federation server implementation allowing for custom handlers per Federation RPC. The example demonstrates configuration, server setup, and interaction with gRPC services.

## Sample Configuration

The Simple Interconnect Provider is configured using a [ConnectorParams](../../pkg/go/config/config.proto) message defined in Protocol Buffers. Below is an example of a configuration in text protobuf format:

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
  - `provider_endpoint_uri`: Specifies the RESTful provider endpoint that the Simple Interconnect Provider bridges.

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
bazel build //examples/golang/simpleinterconnectprovider:simpleinterconnectprovider
```

### Running the Example

Execute the built binary with the necessary flags:

```bash
bazel run //examples/golang/simpleinterconnectprovider:simpleinterconnectprovider -- \
  -config path/to/config.textproto \
  -log-level info
```

#### Flags:

- `-config`: Path to the text protobuf configuration file (e.g., `config.textproto`). Required.
- `-log-level`: Sets the logging level (options: `disabled`, `warn`, `panic`, `info` (default), `fatal`, `error`, `debug`, `trace`).
- `-dry-run`: (Optional) Validate the configuration without starting the server. Exits with a non-zero return code if the config is invalid.

**Example: Dry Run Configuration Validation**

```bash
bazel run //examples/golang/simpleinterconnectprovider:simpleinterconnectprovider -- \
  -config path/to/config.textproto -dry-run
```

For more details, refer to the [main.go](./main.go) source file.

## Sample gRPC Calls

Once the Simple Interconnect Provider is running, you can interact with its gRPC services using `grpcurl`. Below are sample commands for various operations.

### ListTargets

Retrieve service options based on specific filters.

```bash
grpcurl -plaintext -d '{}' localhost:50052 outernet.federation.interconnect.v1alpha.InterconnectService/ListTargets
```

*Refer to [grpc_server.go](../../pkg/go/server/grpc_server.go) and [handler.go](./handler/handler.go) for implementation details.*

### CreateTransceiver

Create a transceiver to connect to the provider's network.

```bash
grpcurl -plaintext -d '{
  "transceiver_id": "my_custom_transceiver",
  "transceiver": {
    "transmit_signal_chain": {
			"antenna": {
				"type": "OPTICAL",
			},
		},    
    "receive_signal_chain": {
			"antenna": {
				"type": "OPTICAL",
			},
		},
		ReceiveSignalChain: &pb.ReceiveSignalChain{
			Antenna: &physical.Antenna{
				Type: physical.Antenna_OPTICAL,
			},
		},
  }
}' localhost:50052 outernet.federation.interconnect.v1alpha.InterconnectService/CreateTransceiver
```

*See [handler.go](./handler/handler.go) for the implementation of `CreateTransceiver`.*

### ListContactWindows

Get the contact windows, where connection between the provider's network and the client's transceiver is possible.

```bash
grpcurl -plaintext -d '{}' localhost:50052 outernet.federation.interconnect.v1alpha.InterconnectService/ListContactWindows
```

*See [handler.go](./handler/handler.go) for the implementation of `ListContactWindows`.*

### DeleteTransceiver

Delete the created transceiver.

```bash
grpcurl -plaintext -d '{ "name": "transceivers/my_custom_transceiver" }' localhost:50052 outernet.federation.interconnect.v1alpha.InterconnectService/DeleteTransceiver
```

*See [handler.go](./handler/handler.go) for the implementation of `DeleteTransceiver`.*

## Project Structure

Understanding the project structure helps in navigating and customizing the Cosmic Connector.

```
/simpleinterconnectprovider/
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
