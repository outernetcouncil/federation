# Glossary of Terms

This section provides definitions for key terms and concepts used throughout the Federation Architecture Specification.

# Key Concepts

## Requestor and Provider Roles

In the Federation API ecosystem, network entities interact primarily through two roles:

- **Requestor**: An entity seeking network services or resources from federated partners. Requestors use the API to discover available resources, request services, and manage ongoing connections.

- **Provider**: An entity offering network services or resources to federated partners. Providers use the API to advertise their capabilities, respond to service requests, and manage resource allocation.

It is entirely possible that bidirectional Federation Services are being offered between two entities, each offering their own Federation API Server for their partner to use as a Requestor.  This flexibility allows for dynamic, multi-directional resource sharing across the federated network.

## Interconnection Points

Interconnection Points are fundamental elements in the Federation API, representing the interfaces where different networks can connect and exchange traffic. They are essential for establishing the physical or logical links that enable federated services.

Characteristics of Interconnection Points:
- Can be physical (e.g., a satellite ground station) or logical (e.g., a virtual network interface)
- Contain information about their capabilities, location, and availability
- May have temporal aspects, especially for mobile or orbiting assets

The `InterconnectionPoint` message in the API provides detailed information about these points, including:
- Unique identifier
- Physical and logical attributes
- Temporal availability
- Associated network interfaces

Understanding and effectively managing Interconnection Points is crucial for creating robust, flexible federated network solutions.

## Service Options

Service Options represent potential service offerings that a Provider can offer to Requestors. They encapsulate the details of a possible network service, allowing Requestors to evaluate and select the most suitable options for their needs.

Components of a Service Option typically include:
- Endpoints (Interconnection Points involved)
- Service attributes (bandwidth, latency, availability)
- Pricing and cost information
- Temporal considerations (time windows for service availability)

The API uses the `ServiceOption` message to convey this information. Requestors can retrieve and analyze Service Options to make informed decisions about which services to request.

Key aspects of working with Service Options:
- Providers generate and update Service Options based on their current network state and policies
- Requestors can filter and sort Service Options based on their requirements
- Selected Service Options form the basis for actual service requests

## General Terms

- **Access Control**: The selective restriction of access to network resources. In Federation, access control mechanisms must span multiple network domains.
- **API**: Application Programming Interface
- **Authentication**: The process of verifying the identity of a user or system. In Federation, robust authentication is crucial for secure inter-network communications.
- **Authorization**: The process of granting or denying access rights to resources. Federation requires careful management of authorization to protect sensitive network resources.
- **Availability**: The time periods during which a network resource or service is accessible for use.
- **Bandwidth**: In the context of computer networking, the maximum rate of data transfer across a given path. In Federation, bandwidth may vary across different network segments and interconnections.
- **Beam Forming**: The ability to dynamically adjust the direction and shape of a satellite's transmission beam to optimize coverage and signal strength for specific areas or users.
- **Bent Pipe Payload**: A satellite transponder that receives, amplifies, and retransmits signals without processing the signal content, essentially acting as a relay in space.
- **Centralized Orchestration**: A model where a single entity manages and coordinates the entire federated network. This approach can provide comprehensive optimization but may face scalability challenges.
- **Distributed Orchestration**: A model where network management is distributed among multiple entities in the Federation. This approach can offer greater scalability and resilience but may result in suboptimal global resource allocation.
- **DoDIN**: Department of Defense Information Network
- **Dynamic Pricing**: A pricing model that adjusts based on real-time demand and availability of network resources. This can help optimize resource utilization across the federated network.
- **Edge Computing**: The practice of processing data near the edge of the network, where it is generated, rather than in a centralized data-processing warehouse. In satellite networks, this can involve processing capabilities on satellites or ground stations.
- **End-to-End Encryption**: A system of communication where only the communicating users can read the messages. This is important for securing data as it traverses multiple network segments in a federated system.
- **Federation API Server**: The server component that implements the Federation API protocol and handles requests from clients, providing the interface between requestors and providers in a federated network.
- **GEO**: Geostationary Earth Orbit
- **Ground Station**: Large, earth-based facilities equipped with substantial antennas and technology to manage communications with satellites and integrate these communications into terrestrial networks.
- **gRPC**: A high-performance, open-source universal RPC framework used as the foundation for the Federation API.
- **Handover**: The process of transferring an active network connection from one satellite or ground station to another as satellites move in their orbits or as user terminals change position.
- **HAPS (High-Altitude Platform Station)**: A telecommunications platform operating in the stratosphere at altitudes of 20-50 km, providing a middle layer of connectivity between terrestrial and satellite networks.
- **HTTP/2**: The underlying transport protocol used by gRPC in the Federation API, providing features such as multiplexing, server push, and header compression.
- **Interconnection Candidate**: One or more interconnection points used to inform the provider of possible interconnections within line-of-sight.
- **Interconnection Point**: A physical or logical point where two federated networks can connect and exchange traffic. These points are crucial for establishing links between different network segments.
- **Interconnection**: An interconnection is a physical and/or logical link through which two networks connect and exchange traffic, enabling communication between requestor and provider networks.
- **Inter-Satellite Link (ISL)**: Communication link between satellites, allowing data to be relayed across a satellite constellation without passing through ground stations.
- **JWT (JSON Web Token)**: A compact, URL-safe means of representing claims to be transferred between two parties, used for authentication and authorization in the Federation API.
- **ISL**: Inter-Satellite Link
- **Latency**: The time delay between the transmission and reception of data.
- **LEO**: Low Earth Orbit
- **Link Budget**: A calculation of all the signal gains and losses in a transmission system, used to determine the quality and feasibility of a communication link.
- **MEO**: Medium Earth Orbit
- **MTU**: Maximum Transmission Unit
- **Network Reachability**: The set of destinations or network prefixes that can be accessed through a given network resource or service. Reachability information is crucial for effective network planning and service provisioning in a federated environment.
- **Network Segment**: A distinct part of a network, such as space, land, or air networks. Each segment may have unique characteristics and requirements.
- **Network Slicing**: A network architecture that enables the creation of multiple virtual networks on top of a shared physical infrastructure, each optimized for specific service requirements in the federated network.
- **Optical Inter-Satellite Link (OISL)**: A high-speed communication link between satellites using laser technology. OISLs enable efficient data transfer within satellite constellations.
- **Protocol Buffers**: A language-agnostic data serialization format used by gRPC for efficient and structured data exchange in the Federation API.
- **Provider**: An entity offering network services or resources to federated partners. Providers use the Federation API to advertise their capabilities, respond to service requests, and manage resource allocation.
- **QoS**: Quality of Service
- **RBAC (Role-Based Access Control)**: A method of regulating access to network resources based on the roles of individual users within the federated network ecosystem.
- **Rate Limiting**: Controls implemented to restrict the number of API requests that can be made within a specified time period, ensuring fair usage and system stability.
- **Reachability**: The set of destinations or network prefixes that can be accessed through a given network resource or service.
- **Requestor**: An entity seeking network services or resources from federated partners. Requestors use the Federation API to discover available resources, request services, and manage ongoing connections.
- **Resource Allocation**: The process of assigning network resources to meet service requirements. In Federation, this process spans multiple network segments and providers.
- **RPC (Remote Procedure Call)**: A protocol that enables a program to execute a procedure in a different address space, forming the basis of the Federation API's communication model.
- **Satellite Constellation**: A group of artificial satellites working together as a system. Constellations can be in various orbits (LEO, MEO, GEO) and serve different purposes in the federated network.
- **SDA**: Space Development Agency
- **Service Level Agreement (SLA)**: A commitment between a service provider and a client, defining the level of service expected. In Federation, SLAs play a crucial role in ensuring quality and reliability across different network segments.
- **Service Option**: A potential service offering from a provider, including details on availability, performance, and cost. Service options allow requestors to evaluate and select suitable network services.
- **Service Request**: A formal request for a specific service, chosen from the advertised service options. It contains details about the desired service, including performance requirements and preferences.
- **Service Status**: Real-time information about an active service in the Federation system, including performance metrics, health indicators, and operational state.
- **Service**: The result of the Federation negotiation. It represents an agreement between requestor and provider for a specific over specific intervals of time.
- **SLA**: Service Level Agreement
- **Spatio-temporal Asset Utilization**: The management and scheduling of network resources considering both their spatial location and temporal availability, particularly important for mobile assets like satellites.
- **Streaming RPC**: A type of RPC that allows for long-lived, bidirectional data streams between clients and servers. Used in the Federation API for real-time updates and continuous data exchange.
- **Temporal Service Attributes**: Time-varying characteristics of a service, such as bandwidth, latency, and availability, which may change over specified time intervals based on network conditions and resource availability.
- **Terrestrial Network**: Ground-based communication infrastructure, including fiber optic cables, cellular towers, and other land-based communication systems.
- **TLS (Transport Layer Security)**: A cryptographic protocol designed to provide communications security over a computer network, used to secure all Federation API communications.
- **Topology**: The arrangement of the physical and logical connections between nodes
- **Unary RPC**: A single request-response style RPC, used in the Federation API for simpler, one-off interactions.
- **User Terminal (UT)**: Devices used by end-users to access satellite network services. They typically include smaller, ground-based antennas designed to communicate directly with satellites.
