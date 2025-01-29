# Federation Go Packages

Core implementation packages for building Federation services in Go.

## Package Structure

```
pkg/go/
├── auth/          # Authentication and authorization
├── interconnectprovider/  # Core Federation Interconnect service implementation
├── handler/       # Federation service interfaces
└── server/        # Server implementations
```

## Core Components

### Authentication (`auth/`)
Provides JWT-based authentication for gRPC services:
- Server interceptors for unary and streaming RPCs
- JWT validation and verification
- RSA public/private key pair support

### Interconnect Provider (`interconnectprovider/`)
Core implementation of the Interconnect service:
- Service lifecycle management
- Configuration handling
- Server component coordination

### Handler Framework (`handler/`)
Interface definitions for implementing Interconnect Federation services:
- `InterconnectHandler` interface
- Support for service scheduling
- Monitoring capabilities
- Service cancellation

### Server Components (`server/`)
Complete server implementations:
- gRPC server for Interconnect API
- Channelz server for monitoring/introspection
- pprof server for profiling
- Generic `Server` interface
- Graceful shutdown support

## Building with Bazel

Common Bazel commands:

```bash
# Build all packages
bazel build //pkg/go/...

# Run all tests
bazel test //pkg/go/...

# Build specific component
bazel build //pkg/go/interconnectprovider
```

## Usage

To use these packages in your own Federation service:

1. Import the required packages:
```go
import (
    "github.com/outernetcouncil/federation/pkg/go/auth"
    "github.com/outernetcouncil/federation/pkg/go/interconnectprovider"
    "github.com/outernetcouncil/federation/pkg/go/handler"
    "github.com/outernetcouncil/federation/pkg/go/server"
)
```

2. Implement the `InterconnectHandler` interface
3. Configure authentication as needed
4. Create and start server components

For a complete implementation example, see the [examples/golang/simpleinterconnectprovider](../../examples/golang/simpleinterconnectprovider) directory.

## References

- [Federation API Reference](../../docs/API_REFERENCE.md)
- [Example Implementation](../../examples/golang/simpleinterconnectprovider)
- [API Reference](https://pkg.go.dev/github.com/outernetcouncil/federation/pkg/go)
