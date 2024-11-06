# Federation Go Packages

Core implementation packages for building Federation services in Go.

## Package Structure

```
pkg/go/
├── auth/          # Authentication and authorization
├── cosmicconnector/  # Core Federation service implementation
├── handler/       # Federation service interfaces
└── server/        # Server implementations
```

## Core Components

### Authentication (`auth/`)
Provides JWT-based authentication for gRPC services:
- Server interceptors for unary and streaming RPCs
- JWT validation and verification
- RSA public/private key pair support

### Cosmic Connector (`cosmicconnector/`)
Core implementation of the Federation service:
- Service lifecycle management
- Configuration handling
- Server component coordination

### Handler Framework (`handler/`)
Interface definitions for implementing Federation services:
- `FederationHandler` interface
- Support for service scheduling
- Monitoring capabilities
- Service cancellation

### Server Components (`server/`)
Complete server implementations:
- gRPC server for Federation API
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
bazel build //pkg/go/cosmicconnector
```

## Usage

To use these packages in your own Federation service:

1. Import the required packages:
```go
import (
    "github.com/outernetcouncil/federation/pkg/go/auth"
    "github.com/outernetcouncil/federation/pkg/go/cosmicconnector"
    "github.com/outernetcouncil/federation/pkg/go/handler"
    "github.com/outernetcouncil/federation/pkg/go/server"
)
```

2. Implement the `FederationHandler` interface
3. Configure authentication as needed
4. Create and start server components

For a complete implementation example, see the [examples/golang/cosmicconnector](../../examples/golang/cosmicconnector) directory.

## References

- [Federation API Guide](../../APIGUIDE.md)
- [Example Implementation](../../examples/golang/cosmicconnector)
- [API Reference](https://pkg.go.dev/github.com/outernetcouncil/federation/pkg/go)
