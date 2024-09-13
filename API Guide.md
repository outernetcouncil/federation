Certainly! Here's a detailed guide to the Outernet Council Federation API:

# Outernet Council Federation API Guide

## Table of Contents

1. Introduction
2. Core Concepts
3. API Overview
4. Detailed API Reference
5. Authentication and Security
6. Best Practices
7. Error Handling
8. Example Workflows
9. Testing and Validation
10. Versioning and Backwards Compatibility

## 1. Introduction

The Outernet Council Federation API enables dynamic orchestration and resource sharing between multiple satellite communication networks. It provides a standardized interface for network operators to advertise resources, discover connectivity options, and provision services across federated networks.

## 2. Core Concepts

### Interconnection Point
An edge node or interface available to form connections between networks. This could be a satellite, ground station, or other network element capable of interfacing with external networks.

### Service Option
A potential service offering from a provider, including endpoints, attributes, and pricing. Service options represent the available connectivity choices that can be requested.

### Service
An active network connection spanning federated resources. Once a service option is selected and provisioned, it becomes an active service.

### Requestor
A network operator requesting services from other networks.

### Provider
A network operator offering services to other networks.

## 3. API Overview

The Federation API consists of several key RPCs:

- StreamInterconnectionPoints: Advertise available interconnection points
- ListServiceOptions: Discover available service offerings
- ScheduleService: Request a specific network service
- MonitorServices: Monitor the status of active services
- CancelService: Terminate an active service

## 4. Detailed API Reference

### StreamInterconnectionPoints

```protobuf
rpc StreamInterconnectionPoints(StreamInterconnectionPointsRequest) 
    returns (stream StreamInterconnectionPointsResponseChunk) {}
```

Allows providers to advertise available interconnection points to requestors.

#### Request
- `snapshot_only` (optional bool): If true, only send a single snapshot of current interconnection points.

#### Response
- `snapshot_complete` (optional SnapshotComplete): Indicates when the initial state snapshot is complete.
- `mutations` (repeated InterconnectionMutation): Updates to advertised interconnection points.

### ListServiceOptions

```protobuf
rpc ListServiceOptions(ListServiceOptionsRequest)
    returns (stream ListServiceOptionsResponse) {}
```

Allows requestors to discover available service options from providers.

#### Request
- `type` (oneof): Specifies the type of service being requested (e.g., requestor_edge_to_requestor_edge, provider_edge_to_provider_edge, provider_edge_to_ip_network).
- `service_attributes_filters` (ServiceAttributesFilterSet): Filters for desired service attributes.

#### Response
- `service_options` (repeated ServiceOption): Available service options matching the request criteria.

### ScheduleService

```protobuf
rpc ScheduleService(ScheduleServiceRequest) 
    returns (ScheduleServiceResponse) {}
```

Allows requestors to request a specific network service from a provider.

#### Request
- `type` (oneof): Specifies the type of service being requested.
- `service_attributes_filters` (ServiceAttributesFilterSet): Desired service attributes.
- `priority` (uint32): Priority of the service request.

#### Response
- `service_id` (string): Unique identifier for the scheduled service.

### MonitorServices

```protobuf
rpc MonitorServices(stream MonitorServicesRequest)
    returns (stream MonitorServicesResponse) {}
```

Allows requestors to monitor the status of active services.

#### Request
- `add_service_ids` (repeated string): Services to start monitoring.
- `drop_service_ids` (repeated string): Services to stop monitoring.

#### Response
- `updated_services` (repeated ServiceStatus): Status updates for monitored services.

### CancelService

```protobuf
rpc CancelService(CancelServiceRequest) 
    returns (CancelServiceResponse) {}
```

Allows requestors to cancel an active service.

#### Request
- `service_id` (string): Identifier of the service to cancel.

#### Response
- `cancelled` (bool): Indicates if the service was successfully cancelled.

## 5. Authentication and Security

The Federation API uses OAuth 2.0 for authentication. Clients must obtain an access token and include it in the `Authorization` header of each request:

```
Authorization: Bearer <access_token>
```

All API communications should be encrypted using TLS 1.2 or later.

## 6. Best Practices

- Keep interconnection point information up-to-date by frequently calling StreamInterconnectionPoints.
- Use appropriate filters in ListServiceOptions to reduce unnecessary data transfer.
- Implement exponential backoff for retrying failed requests.
- Monitor active services regularly to detect and respond to changes quickly.

## 7. Error Handling

The API uses standard gRPC status codes. Common errors include:

- INVALID_ARGUMENT (3): Request parameters are invalid
- NOT_FOUND (5): Requested resource doesn't exist
- ALREADY_EXISTS (6): Attempt to create a resource that already exists
- PERMISSION_DENIED (7): Lack of necessary permissions
- RESOURCE_EXHAUSTED (8): Quota or rate limit exceeded

Detailed error messages are provided in the `details` field of the status response.

## 8. Example Workflows

### Requesting a Service

1. Call StreamInterconnectionPoints to discover available interconnection points.
2. Use ListServiceOptions to find suitable service offerings.
3. Call ScheduleService with the chosen service option.
4. Monitor the service status using MonitorServices.
5. When no longer needed, call CancelService.

## 9. Testing and Validation

[TODO]

## 10. Versioning and Backwards Compatibility

[TODO]
