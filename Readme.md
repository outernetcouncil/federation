# Outernet Council Federation API

The Outernet Council Federation API enables dynamic orchestration and resource sharing between multiple satellite communication networks.

## Overview

The Federation API allows network operators to:

- Advertise available network resources and interconnection points
- Discover resources and connectivity options across federated networks  
- Request and reserve network services spanning multiple providers
- Coordinate spectrum usage and avoid interference
- Dynamically provision end-to-end connectivity leveraging federated resources

Key features:

- Open, vendor-agnostic API for space network federation
- Support for LEO, MEO, GEO satellite constellations and ground networks
- Dynamic resource availability advertising and discovery
- Service request/offer negotiation and reservation  
- Real-time network orchestration across federated segments

## API Components

The Federation API consists of the following key components:

### StreamInterconnectionPoints

Allows providers to advertise available interconnection points (network edges) to requestors.

```protobuf
rpc StreamInterconnectionPoints(StreamInterconnectionPointsRequest) 
    returns (stream StreamInterconnectionPointsResponseChunk) {}
```

### ListServiceOptions  

Allows requestors to discover available service options from providers.

```protobuf
rpc ListServiceOptions(ListServiceOptionsRequest)
    returns (stream ListServiceOptionsResponse) {}
```

### ScheduleService

Allows requestors to request a specific network service from a provider.

```protobuf
rpc ScheduleService(ScheduleServiceRequest) 
    returns (ScheduleServiceResponse) {}
```

### MonitorServices

Allows requestors to monitor the status of active services.

```protobuf
rpc MonitorServices(stream MonitorServicesRequest)
    returns (stream MonitorServicesResponse) {}
```

### CancelService

Allows requestors to cancel an active service.

```protobuf
rpc CancelService(CancelServiceRequest) 
    returns (CancelServiceResponse) {}
```

## Key Concepts

- **Interconnection Point**: An edge node or interface available to form connections between networks.

- **Service Option**: A potential service offering from a provider, including endpoints, attributes, and pricing.

- **Service**: An active network connection spanning federated resources.

- **Requestor**: Network operator requesting services from other networks.

- **Provider**: Network operator offering services to other networks.

## Usage Flow

1. Providers advertise available interconnection points via `StreamInterconnectionPoints`

2. Requestors discover service options via `ListServiceOptions`

3. Requestors select and request services via `ScheduleService` 

4. Active services are monitored via `MonitorServices`

5. Services can be cancelled via `CancelService` when no longer needed

## Detailed Guide

[Coming Soon]

## Documentation  

[Coming Soon]

## Contributing

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
