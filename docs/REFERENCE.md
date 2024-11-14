# Federation Reference Architecture v0.1

# TOC

- [TOC](#toc)  
- [1\. Introduction](#1-introduction)  
  - [1.1. Purpose of the Federation Architecture Specification](#11-purpose-of-the-federation-architecture-specification)  
  - [1.2. Scope of the Document](#12-scope-of-the-document)  
  - [1.3. Background and Context](#13-background-and-context)  
  - [1.4. Key Use Cases](#14-key-use-cases)  
    - [Multi-Orbit Satellite Network Integration](#multi-orbit-satellite-network-integration)  
    - [Military and Defense Network Commercial Interoperability](#military-and-defense-network-commercial-interoperability)  
    - [Commercial Aviation Connectivity](#commercial-aviation-connectivity)  
    - [Optical Ground Station as a Service (OGSaaS)](#optical-ground-station-as-a-service-ogsaas)  
    - [Disaster Response and Emergency Communications](#disaster-response-and-emergency-communications)  
  - [1.5 Intended Audience](#15-intended-audience)  
  - [1.6. Related Documents and Resources](#16-related-documents-and-resources)  
  - [1.7. Versioning and Updates](#17-versioning-and-updates)  
- [2\. Terms and Definitions](#2-terms-and-definitions)  
  - [2.1. General Terms](#21-general-terms)  
  - [2.2. Federation-Specific Terminology](#22-federation-specific-terminology)  
  - [2.3. Network Types and Components](#23-network-types-and-components)  
  - [2.4. Technical Concepts](#24-technical-concepts)  
  - [2.5. API and Protocol Terms](#25-api-and-protocol-terms)  
  - [2.6. Orchestration and Management Concepts](#26-orchestration-and-management-concepts)  
  - [2.7. Security and Privacy Terms](#27-security-and-privacy-terms)  
  - [2.8. Acronyms](#28-acronyms)  
- [3\. Architecture Description Overview](#3-architecture-description-overview)  
  - [3.1. Federation Architecture Vision](#31-federation-architecture-vision)  
  - [3.2. Key Principles and Concepts](#32-key-principles-and-concepts)  
  - [3.3. Federation Models](#33-federation-models)  
    - [3.3.1. Peer-to-Peer Model](#331-peer-to-peer-model)  
    - [3.3.2. Centralized (Hub) Model](#332-centralized-hub-model)  
  - [3.4. API Structure and Core Resources](#34-api-structure-and-core-resources)  
  - [3.5. Federation Workflow Overview](#35-federation-workflow-overview)  
  - [3.6. Temporal Considerations](#36-temporal-considerations)  
  - [3.7. Security and Privacy Considerations](#37-security-and-privacy-considerations)  
  - [3.8. Scalability and Performance](#38-scalability-and-performance)  
  - [3.9. Interoperability and Standards](#39-interoperability-and-standards)  
  - [3.10. Future Extensibility](#310-future-extensibility)  
- [4\. Architecture Rationale](#4-architecture-rationale)  
  - [4.1. Key Architectural Decisions](#41-key-architectural-decisions)  
    - [4.1.1. Support for Multiple Federation Models](#411-support-for-multiple-federation-models)  
    - [4.1.2. gRPC and Protocol Buffers for API Implementation](#412-grpc-and-protocol-buffers-for-api-implementation)  
    - [4.1.3. Temporal Considerations in API Objects](#413-temporal-considerations-in-api-objects)  
    - [4.1.4. Dynamic Service Option Generation](#414-dynamic-service-option-generation)  
  - [4.2. Trade-offs and Alternatives](#42-trade-offs-and-alternatives)  
    - [4.2.1. Peer-to-Peer vs. Centralized Orchestration](#421-peer-to-peer-vs-centralized-orchestration)  
    - [4.2.2. Service Option Generation Approaches](#422-service-option-generation-approaches)  
    - [4.2.3. Temporal Granularity in API Objects](#423-temporal-granularity-in-api-objects)  
  - [4.3. Federation Models](#43-federation-models)  
    - [4.3.1. Peer-to-Peer (Requestor-Provider) Model](#431-peer-to-peer-requestor-provider-model)  
    - [4.3.2. Centralized (Hub) Model](#432-centralized-hub-model)  
  - [4.4. Service Option Approaches](#44-service-option-approaches)  
    - [4.4.1. Provider-Only Path](#441-provider-only-path)  
    - [4.4.2. Provider \+ Requestor Resource Path](#442-provider--requestor-resource-path)  
    - [4.4.3. Hybrid and Adaptive Approaches](#443-hybrid-and-adaptive-approaches)  
    - [4.4.4. Considerations for Implementation](#444-considerations-for-implementation)  
  - [4.5. Temporal Considerations](#45-temporal-considerations)  
  - [4.6. Security and Privacy Design Decisions](#46-security-and-privacy-design-decisions)  
  - [4.7. Scalability and Performance Considerations](#47-scalability-and-performance-considerations)  
  - [4.8. Interoperability Decisions](#48-interoperability-decisions)  
  - [4.9. Future Extensibility Considerations](#49-future-extensibility-considerations)

# 1\. Introduction

## 1.1. Introduction to Federation

## Federation is a protocol that connects networks together in a way that unlocks new capabilities and enables seamless interoperability between participant networks. Bringing together various interoperable networks with their unique strengths and capabilities enables more complex and capable user experiences. As an open community led by its contributors, Federation’s architecture is evolving as network technologies are brought in to adapt to evolving customer needs. The Outernet Council’s development and evolution of Federation is guided by how the network can serve the user’s needs along the following themes:

* Reducing costs  
* Improving choice  
* Network resilience  
* Network security  
* Enabling new capabilities

Federation is part of a larger vision of a connected, secure, resilient, network of communications satellites provisioned by a diverse ecosystem of operators across multiple countries – the Outernet. The Outernet has the potential to unify space and terrestrial network architectures across diverse network segments, including land, sea, air and space.

## 1.2 Purpose of the Federation Architecture Specification

The primary purposes of this specification are to:

1. Establish a standardized approach for integrating diverse network resources  
2. Enable dynamic allocation of resources across federated networks  
3. Facilitate seamless communication between terrestrial, aerial, and space-based systems  
4. Provide a framework for expanding market reach and optimizing resource utilization  
5. Accelerate innovation in global connectivity solutions

Additionally, this specification directly addresses the complexities of spatio-temporal asset utilization, such as scheduling assets with complex time-dependent availability and capacity.

## 1.2. Scope of the Document

This Federation Architecture Specification encompasses:

- The definition of the Federation Architecture  
- Detailed coverage of systems, networks, and technologies involved in Federation  
- Integration strategies for land, sea ,air and space network segments  
- Specifications for both discrete physical resource sharing and spectrum sharing  
- API definitions and protocols for inter-network communication  
- Implementation guidelines for network operators and service providers

The document addresses the technical and operational aspects of network Federation, providing a comprehensive guide for stakeholders involved in the development and deployment of federated network solutions.

## 1.3. Background and Context

The Federation Architecture Specification emerges in response to the rapidly evolving landscape of global communications. As the boundaries between terrestrial and space-based networks continue to blur, there is an increasing need for a unified approach to network integration and resource sharing.

Historically, different network segments \- terrestrial, satellite, and aerial \- have operated largely in isolation, with limited interoperability. This siloed approach has led to inefficiencies in resource utilization, gaps in global coverage, and challenges in providing seamless connectivity across diverse environments.

Recent technological advancements, particularly in satellite communications with the proliferation of Low Earth Orbit (LEO) constellations, have opened new possibilities for global connectivity. Simultaneously, the demand for ubiquitous, high-bandwidth communications has surged across various sectors, including commercial, military, and emergency services.

The concept of network Federation has gained traction as a solution to these challenges. Federation allows diverse network operators to share resources, extend their reach, and optimize service delivery. Past attempts at spectrum and/or capacity interoperability have demonstrated its potential but also highlighted the need for a standardized[^1], comprehensive approach.

This specification builds upon these early efforts and lessons learned from initial implementations of federated systems. It aims to address key challenges such as:

1. Interoperability between disparate network technologies and protocols  
2. Dynamic resource allocation across multiple network domains  
3. Security and privacy concerns in shared network environments  
4. Complex spatio-temporal considerations, especially in satellite-based systems  
5. Scalability to accommodate growing networks and increasing data demands

By providing a robust framework for Federation, this specification seeks to enable a new era of global connectivity, where terrestrial, satellite, and aerial networks seamlessly integrate to provide ubiquitous, efficient, and resilient communication services.

## 1.4. Key Use Cases

The Federation Architecture is designed to support a wide range of use cases, including but not limited to:

### Space-Relay Interconnect and Transit

* Government satellite interconnecting in space (via RF or optical) with a commercial constellation for on-demand network transit to a given destination

### Ground Segment as a Service

* Dynamic reservation of third party ground stations (optical or RF) for on-demand network transit to specified points of presence, ops centers, or storage

### “Agile” (Multi-Constellation) User Terminals

* Multiple projects are underway to develop "agile" user terminals with multiple physical or virtualized modems that share a single aperture and are capable of dynamically switching between multiple commercial or military SATCOM systems  
* Federation enables just-in time on-demand ordering and provisioning of beams and transit for fixed terrestrial, land mobile, airborne, or maritime terminals to roam between these networks

### Spectrum Coordination \[new\]

\- TODO

### Service Requests

#### Optimization Requests

Examples: 

* **EO operator wants to land traffic as quickly as possible**  
  * Using the federation API, the operator requests providers (satellites or groundstations) that can make contact and land the data. The request allows specification of desired spectrum, required SLA, size of data, security requirements, etc.  
  * The federation engine queries available providers and generates a sorted list of options and their corresponding costs.   
  * Requesting operator selects an option to the federation engine which processes payment and issues network instructions to the involved parties  
  * Federation engine monitors the network KPIs and generates alerts   
* **EO operator wants to find the cheapest way to get data down**  
  * Using the federation API, the operator requests providers (satellites or groundstations) that can make contact and land the data. The request allows specification of desired spectrum, required SLA, size of data, security requirements, etc.  
  * The federation engine queries available providers and generates a sorted list of options and their corresponding costs.   
  * Requesting operator selects an option to the federation engine which processes payment and issues network instructions to the involved parties  
  * Federation engine monitors the network KPIs and generates alerts

#### Multi party requests

Examples: 

* **Newer LEO provider wants to start commercial service before MVP constellation has been built out**  
  * For the gaps in coverage for the newly launching constellation, the provider requests coverage via the federation API. The request allows specification of desired spectrum, required SLA, expected number of subscribers, security requirements, etc.  
  *  The federation engine provides a list of willing providers (LEO, MEO or GEO) for the provider to choose.  
  * This process can be dynamic as the LEO provider launches more satellites, onboards various customers or seeks to provide specific SLAs.  
  * Requesting operator selects an option to the federation engine which processes payment and issues network instructions to the involved parties  
  * Federation engine monitors the network KPIs and generates alert  
* **Fire monitoring service wants to source real time data**  
  * Application provider uses the Federation API to request all available data streams (LEO Earth observation satellite, GEO radar satellite, HAPS thermal monitoring, Optical systems of helicopters, etc) for a given area to be routed to their data lake in real time  
  * Federation engine provides the network instructions and the API keys needed to route traffic to the desired destination  
  * Federation engine monitors the network KPIs and issues alerts  
  * For TS-SDN tenants using the east-west interface, the federation API enables coordination between networks as they anticipate and preempt disruptions

### Resource sharing

#### Temporary Allocation

Examples: 

* **LEO operator needs to take a ground station offline for maintenance.**  
  * Using the federation API, the operator requests alternate ground stations to land data. The request allows specification of desired spectrum, duration of request, security requirements, etc.  
  * The federation engine queries available providers and generates a sorted list of options and their corresponding costs.   
  * Requesting operator selects an option to the federation engine which processes payment and issues network instructions to the involved parties  
  * Federation engine monitors the network KPIs and generates alerts   
* **LEO operator sees higher than serviceable demand for a region**  
  * Using the federation API, the operator requests alternate constellations for available front haul spectrum. The request allows specification of desired spectrum, region, duration of request, security requirements, etc.  
  * The federation engine queries available providers and generates a sorted list of options and their corresponding costs.   
  * Requesting operator selects an option to the federation engine which processes payment and issues network instructions to the involved parties  
  * Federation engine monitors the network KPIs and generates alerts 

## 1.5 Intended Audience

This specification is intended for:

- System architects of existing network operators  
- Network operators and service providers  
- Satellite communication system developers  
- Terrestrial network engineers  
- Software developers working on network integration  
- Researchers and academics in the field of global communications  
- Policymakers and regulators involved in spectrum management and network interoperability

Readers are expected to have a basic understanding of network architectures, satellite communications, and API concepts. 

## 1.6. Related Documents and Resources

For a comprehensive understanding of Federation, the following related documents and resources may be consulted:

- Existing documentation from [Outernet Council's Federation Repo](https://github.com/outernetcouncil/federation)  
  - [APIGUIDE.md](http://APIGUIDE.md), a detailed guide to supporting the Federation API  
  - [README.md](http://README.md), a succinct motivation summary of the Federation API  
  - [Federation.proto](http://federation.proto) Protocol Buffer definitions for the Federation Service  
- [Outernet Council's Network Model for Temporospatial Systems (NMTS) Repo](https://github.com/outernetcouncil/nmts)  
  - This comprehensive network model supports network representation across space, ground, air paradigms

Additional resources and example implementations will be made available in the [Outernet's Federation Repo](https://github.com/outernetcouncil/federation).

## 1.7. Versioning and Updates

The Federation Architecture Specification follows semantic versioning principles:

- MAJOR version for incompatible API changes  
- MINOR version for backwards-compatible functionality additions  
- PATCH version for backwards-compatible bug fixes

This document represents a pre-version 1.0.0 alpha of the Federation Architecture Specification.

Updates and revisions to this specification will be managed through a controlled process, with major updates typically released on an annual basis. Minor updates and patches may be released more frequently as needed.

To access the most current version of this specification and stay informed about updates:

1. Subscribe to the Federation Architecture mailing list  
2. Monitor the public GitHub repository for change logs and release notes

Feedback and contributions to the evolution of this specification are welcome from the community and industry partners.

# 2\. Terms and Definitions

This section provides definitions for key terms and concepts used throughout the Federation Architecture Specification. Understanding these terms is crucial for comprehending the architecture and its components.

## 2.1. General Terms

- **Access Control**: The selective restriction of access to network resources. In Federation, access control mechanisms must span multiple network domains.  
- **API**: Application Programming Interface  
- **Authentication**: The process of verifying the identity of a user or system. In Federation, robust authentication is crucial for secure inter-network communications.  
- **Authorization**: The process of granting or denying access rights to resources. Federation requires careful management of authorization to protect sensitive network resources.  
- **Bandwidth**: In the context of computer networking, the maximum rate of data transfer across a given path. In Federation, bandwidth may vary across different network segments and interconnections.  
- **Bent Pipe Payload**: A satellite transponder that receives, amplifies, and retransmits signals without processing the signal content, essentially acting as a relay in space.  
- **Centralized Orchestration**: A model where a single entity manages and coordinates the entire federated network. This approach can provide comprehensive optimization but may face scalability challenges.  
- **Distributed Orchestration**: A model where network management is distributed among multiple entities in the Federation. This approach can offer greater scalability and resilience but may result in suboptimal global resource allocation.  
- **DoDIN**: Department of Defense Information Network  
- **Dynamic Pricing**: A pricing model that adjusts based on real-time demand and availability of network resources. This can help optimize resource utilization across the federated network.  
- **End-to-End Encryption**: A system of communication where only the communicating users can read the messages. This is important for securing data as it traverses multiple network segments in a federated system.  
- **Federation**: The act of combining multiple independent networks to create a larger, more capable network ecosystem. In the context of this architecture, it refers to the seamless integration of diverse network resources, including terrestrial, aerial, and space-based systems.  
- **GEO**: Geostationary Earth Orbit  
- **Ground Station**: Large, earth-based facilities equipped with substantial antennas and technology to manage communications with satellites and integrate these communications into terrestrial networks.  
- **gRPC**: A high-performance, open-source universal RPC framework used as the foundation for the Federation API.  
- **Handover**: The process of transferring an active network connection from one satellite or ground station to another as satellites move in their orbits or as user terminals change position.  
- **Interconnection Candidate**: One or more interconnection points used to inform the provider of possible interconnections within line-of-sight.  
- **Interconnection Point**: A physical or logical point where two federated networks can connect and exchange traffic. These points are crucial for establishing links between different network segments.  
- **Interconnection**: An interconnection is a physical and/or logical link through which two networks connect and exchange traffic, enabling communication between requestor and provider networks.  
- **Inter-Satellite Link (ISL)**: Communication link between satellites, allowing data to be relayed across a satellite constellation without passing through ground stations.  
- **ISL**: Inter-Satellite Link  
- **Latency**: The time delay between the transmission and reception of data.   
- **LEO**: Low Earth Orbit  
- **Link Budget**: A calculation of all the signal gains and losses in a transmission system, used to determine the quality and feasibility of a communication link.  
- **MEO**: Medium Earth Orbit  
- **MTU**: Maximum Transmission Unit  
- **Network Reachability**: The set of destinations or network prefixes that can be accessed through a given network resource or service. Reachability information is crucial for effective network planning and service provisioning in a federated environment.  
- **Network Segment**: A distinct part of a network, such as space, land, or air networks. Each segment may have unique characteristics and requirements.  
- **Optical Inter-Satellite Link (OISL)**: A high-speed communication link between satellites using laser technology. OISLs enable efficient data transfer within satellite constellations.  
- **Protocol Buffers**: A language-agnostic data serialization format used by gRPC for efficient and structured data exchange in the Federation API.  
- **Provider**: An entity offering network services or resources to federated partners. Providers use the Federation API to advertise their capabilities, respond to service requests, and manage resource allocation.  
- **QoS**: Quality of Service  
- **Requestor**: An entity seeking network services or resources from federated partners. Requestors use the Federation API to discover available resources, request services, and manage ongoing connections.  
- **Resource Allocation**: The process of assigning network resources to meet service requirements. In Federation, this process spans multiple network segments and providers.  
- **Satellite Constellation**: A group of artificial satellites working together as a system. Constellations can be in various orbits (LEO, MEO, GEO) and serve different purposes in the federated network.  
- **SDA**: Space Development Agency  
- **Service Level Agreement (SLA)**: A commitment between a service provider and a client, defining the level of service expected. In Federation, SLAs play a crucial role in ensuring quality and reliability across different network segments.  
- **Service Option**: A potential service offering from a provider, including details on availability, performance, and cost. Service options allow requestors to evaluate and select suitable network services.  
- **Service Request**: A formal request for a specific service, chosen from the advertised service options. It contains details about the desired service, including performance requirements and preferences.  
- **Service**: The result of the Federation negotiation. It represents an agreement between requestor and provider for a specific over specific intervals of time.  
- **SLA**: Service Level Agreement  
- **Streaming RPC**: A type of RPC that allows for long-lived, bidirectional data streams between clients and servers. Used in the Federation API for real-time updates and continuous data exchange.  
- **Terrestrial Network**: Ground-based communication infrastructure, including fiber optic cables, cellular towers, and other land-based communication systems.  
- **Topology**: The arrangement of the physical and logical connections between nodes  
- **Unary RPC**: A single request-response style RPC, used in the Federation API for simpler, one-off interactions.  
- **User Terminal (UT)**: Devices used by end-users to access satellite network services. They typically include smaller, ground-based antennas designed to communicate directly with satellites.

# 3\. Architecture Description Overview

## 3.1. Federation Architecture Vision

The Federation Architecture aims to create a unified framework for seamlessly integrating space and terrestrial network architectures. This vision addresses the growing need for global connectivity across diverse network segments, including space, land, and air networks.

Key objectives of the Federation Architecture include:

- Enabling dynamic sharing of information, resources, and spectrum across networks  
- Facilitating seamless communication between terrestrial, aerial, and space-based systems  
- Optimizing resource utilization through intelligent allocation and Federation  
- Expanding market reach for network operators and service providers  
- Accelerating innovation in global connectivity solutions

The architecture is designed to handle the complexities of spatio-temporal asset utilization, allowing for efficient scheduling of assets with complex time-dependent availability and capacity.

## 3.2. Key Principles and Concepts

The Federation Architecture is built upon several fundamental principles and concepts that guide its design and implementation:

1. **Network Agnosticism**: The architecture is designed to accommodate diverse network types, including terrestrial, satellite, and aerial systems, without favoring any specific technology or provider.  
2. **Dynamic Resource Allocation**: The framework enables real-time allocation and reallocation of network resources based on changing demands, availability, and network conditions and capabilities.  
3. **Temporal Awareness**: Recognizing the time-varying nature of network resources, especially in satellite-based systems, where the architecture incorporates temporal considerations in all aspects of service definition and management.  
4. **Scalability**: The design supports scaling from small-scale Federations to large, complex networks involving multiple providers, network segments, and diverse technologies.  
5. **Interoperability**: Standardized interfaces and protocols ensure seamless communication and integration between different network segments and providers.  
6. **Security and Privacy by Design**: The architecture incorporates robust security measures and privacy-preserving techniques as core components rather than afterthoughts.  
7. **Flexibility in Federation Models**: Support for both peer-to-peer and multi-party Federation models allows adaptation to various operational and business requirements.  
8. **Service-Oriented Approach**: The architecture is built around the concept of services, allowing for clear definition, negotiation, and management of network capabilities.  
9. **Autonomy and Control**: While enabling Federation, the architecture respects the autonomy of individual network operators, allowing them to maintain control over their resources and participation levels.  
10. **Efficiency**: The architecture aims to optimize resource utilization across federated networks, maximizing the value of existing infrastructure investments.  
11. **Resilience**: By enabling diverse network integration, the framework enhances overall network resilience, providing multiple pathways for communication and service delivery.  
12. **Global Connectivity Vision**: The ultimate goal of the architecture is to facilitate ubiquitous, high-quality connectivity on a global scale, bridging gaps between traditionally siloed network infrastructures.

These principles and concepts form the foundation of the Federation Architecture, guiding its development and implementation to create a robust, flexible, and future-proof framework for integrating diverse network resources across space and terrestrial domains.

## 3.3. Federation Models

The Federation Architecture supports two primary models for network interaction: the Peer-to-Peer (Requestor-Provider) Model and the Multi-Party (Hub) Model. Each model has its own set of considerations, strengths, and potential drawbacks.

### 3.3.1. Peer-to-Peer Model

In this distributed model, communication occurs directly between requestors and providers without a central intermediary.

**Considerations:**

- Requires direct business agreements between operators  
- May result in increased complexity when interacting with multiple Federation partners  
- Potential for suboptimal global resource allocation

**Strengths:**

1. Enhanced Privacy: Minimizes data sharing with third parties  
2. Greater Control: Operators maintain direct control over their resources and sharing policies  
3. Reduced Single Point of Failure: No central hub that could disrupt all Federation activities if it fails  
4. Flexibility: Allows for customized agreements between specific partners  
5. Scalability: Can grow organically as new partners join the Federation

### 3.3.2. Multi-Party Model

This model introduces a central orchestrator that facilitates the fulfillment of service requests between multiple requestors and providers.

**Considerations:**

- Requires trust in the central entity to manage interactions fairly  
- Requires design for resilience to ensure no single point of failure

**Strengths:**

1. Simplified Operations: Reduces complexity for individual network operators  
2. Global Optimization: Potential for more efficient overall resource allocation  
3. Standardization: Can enforce consistent protocols and data formats across the Federation  
4. Market Efficiency: Facilitates easier discovery of available resources and services  
5. Centralized Monitoring and Management: Enables comprehensive oversight of Federation activities

The Federation Architecture is designed to support multiple forms of centralization, from fully distributed to fully centralized, allowing for flexibility in implementation based on specific use cases and partner requirements. Hybrid models, combining elements of both peer-to-peer and centralized approaches, are also possible within this framework.

Additionally, the first instances of this type of Federation would most likely need to be proven on a P2P basis (e.g. MNO leasing out spectrum directly to one or more partnered SNOs on a primary-secondary basis) where the individual actors assume all responsibility, before a regulator would feel confident to be the party which facilitates the reuse of spectral resources between otherwise contending operators.

## 3.4. API Structure and Core Resources

The Federation API is built on gRPC and uses Protocol Buffers for efficient data serialization. See 1.6. Related Documents and Resources for references to detailed API explanation and source.

## 

## 3.5. Federation Workflow Overview

The following sequence diagram gives the high level lifecycle for a Federation request.  
![][image1]  
*A sample Federation workflow between peers which supports Query (discovery of service opportunities that meet some Requestor needs), Request (explicit ask for services from the Provider), Monitoring (of resources, pricing, service planned vs reported attributes), and finally Termination*.”

## 3.6. Temporal Considerations

The Federation Architecture places significant emphasis on temporal aspects:

- Most API objects (e.g., ServiceOptions, Services) include time intervals to represent their validity periods  
- The architecture handles dynamic changes in network topology, e.g. for satellite-based systems or vehicle mounted terminals  
- Service attributes (bandwidth, latency, availability) are often represented as temporal values

## 3.7. Security and Privacy Considerations

Security and privacy are paramount in the Federation Architecture:

- Authentication and authorization mechanisms are built into the API  
- The architecture allows for varying levels of information sharing based on trust relationships  
- Privacy-preserving techniques are employed, especially in the peer-to-peer model  
- End-to-end encryption is supported for sensitive data transmission

## 3.8. Scalability and Performance

The architecture is designed to handle large-scale Federations:

- Efficient data structures and algorithms for managing time-series data  
- Support for streaming updates to handle real-time changes in network conditions  
- Optimization techniques for service option generation and evaluation

## 3.9. Interoperability and Standards

The Federation Architecture aims to ensure interoperability:

- Alignment with relevant industry standards (e.g., 3GPP, ETSI, ITU recommendations)  
- Support for common network protocols and data formats  
- See [Outernet Council's Network Model for Temporospatial Systems (NMTS)](https://github.com/outernetcouncil/nmts) used by Federation API

## 3.10. Future Extensibility

The architecture is designed with future growth in mind:

- Modular design allowing for the addition of new service types and network technologies  
- Versioning system to manage API evolution while maintaining backwards compatibility  
- Extensible data models to accommodate emerging use cases and requirements

By adhering to these architectural principles and leveraging the power of Federation, the architecture aims to create a flexible, scalable, and efficient framework for global network integration across space and terrestrial domains.

# 4\. Architecture Rationale

This section provides insights into the key architectural decisions, trade-offs, and alternatives considered in the design of the Federation Architecture. It aims to explain the reasoning behind major design choices and how they address the unique challenges of integrating diverse network resources across space and terrestrial domains.

## 4.1. Key Architectural Decisions

### 4.1.1. Support for Multiple Federation Models

**Decision**: The architecture supports both peer-to-peer and centralized Federation models.

**Rationale**: This flexibility allows the architecture to accommodate various business relationships, privacy requirements, and operational preferences of different network operators. It also enables a gradual transition from simpler peer-to-peer interactions to more complex, centralized orchestration as the Federation ecosystem matures.

### 4.1.2. gRPC and Protocol Buffers for API Implementation

**Decision**: The Federation API is implemented using gRPC with Protocol Buffers.

**Rationale**: This combination offers several advantages:

- Efficient binary serialization, reducing network overhead  
- Strong typing enhances code reliability and catches errors at compile time  
- Language-agnostic implementation, facilitating integration across diverse systems  
- Built-in support for streaming RPCs, crucial for real-time updates in dynamic network environments

### 4.1.3. Temporal Considerations in API Objects

**Decision**: Incorporate time intervals and temporal attributes in most API objects.

**Rationale**: Given the dynamic nature of satellite constellations and the time-varying availability of network resources, temporal awareness is crucial. This decision enables accurate representation of:

- Satellite visibility windows, mobile terminals, etc  
- Time-varying service attributes (e.g., bandwidth, latency)  
- Scheduled maintenance or resource availability

### 4.1.4. Dynamic Service Option Generation

**Decision**: Implement on-demand generation of service options rather than pre-computing all possibilities.

**Rationale**: This approach balances computational efficiency with the need for up-to-date and relevant service options. It reduces the storage and processing requirements, especially for large constellations, while ensuring that requestors receive current and applicable options.

## 4.2. Trade-offs and Alternatives

### 4.2.2. Service Option Generation Approaches

**Trade-off**: Computational burden vs. accuracy of service attributes

- **Provider-Only Path**:  
    
  - Pros: Simpler implementation, less information shared by requestor  
  - Cons: Potential inaccuracies in link evaluation, especially for multipoint or coverage beams


- **Provider \+ Requestor Resource Path**:  
    
  - Pros: More accurate link evaluation, better handling of complex beam types  
  - Cons: Increased computational burden on provider, more information shared by requestor

**Decision**: Implement both approaches, allowing selection based on specific scenarios and privacy requirements.

### 4.2.3. Temporal Granularity in API Objects

**Trade-off**: Accuracy vs. computational and network overhead

- Fine-grained intervals provide more accurate representation but increase data volume and processing requirements  
- Coarse-grained intervals reduce overhead but may miss short-term variations in resource availability or performance

**Decision**: Allow flexible temporal granularity, with recommendations for common scenarios to balance accuracy and efficiency.

## 4.3. Federation Models

### 4.3.1. Peer-to-Peer (Requestor-Provider) Model

In this distributed model:

- Communication occurs directly between requestors and providers  
- Offers greater privacy and control for requestors  
- Requires individual network controllers (e.g., instances of Spacetime) to make decisions  
- Necessitates business agreements between operators to determine information sharing protocols

### 4.3.2. Multi-Party Model

This model introduces multiple parties interacting with one or more hubs:

- Facilitates coordination between multiple requestors and providers  
- Potentially offers greater optimization of resource allocation  
- Reduces the complexity for individual network operators  
- May provide standardization of data and processes across the Federation

The architecture is designed to support a spectrum of centralization, from fully distributed to fully centralized, allowing for flexibility in implementation based on specific use cases and partner requirements.

## 4.4. Capacity: Service Option Approaches

The Federation Architecture supports two primary approaches for generating and evaluating service options: the Provider-Only Path and the Provider \+ Requestor Resource Path. Each approach has its own set of advantages, challenges, and use cases.

### 4.4.1. Provider-Only Path

In this approach, service options are generated and evaluated solely based on the Provider's network information.

**Key Characteristics:**

- Service options have Provider interconnection points as endpoints  
- Requires minimal information from the Requestor  
- Provider has full control over the service option generation process

**Advantages:**

1. **Simplicity**: Easier to implement and manage from the Provider's perspective  
2. **Privacy**: Requires less information sharing from the Requestor, enhancing privacy  
3. **Reduced Computational Burden on Requestor**

**Challenges:**

1. **Scalability**: Can face challenges in large constellations (e.g., N^2 problem with satellite interconnections)  
2. **Accuracy Limitations**: May lead to inaccuracies in link evaluation, especially for multipoint or coverage beams  
3. **Limited Optimization**: Without detailed Requestor information, it's harder to optimize for specific Requestor needs

**Use Cases:**

- Initial service discovery phase where Requestors want to explore options without sharing detailed information  
- Scenarios with simple, point-to-point connections where Provider information is sufficient for accurate evaluation  
- When privacy concerns outweigh the need for highly optimized service options

### 4.4.2. Provider \+ Requestor Resource Path

This approach involves both the Provider and Requestor in the service option generation and evaluation process.

**Key Characteristics:**

- Service options include Requestor interconnection points as endpoints  
- Requires more information exchange between Provider and Requestor  
- Allows for more accurate and tailored service options

**Advantages:**

1. **Accuracy**: Enables more precise link evaluation, especially for complex scenarios like multipoint or coverage beams  
2. **Optimization**: Better ability to tailor service options to Requestor's specific needs and constraints  
3. **Flexibility**: Can handle a wider range of complex network configurations and service requirements  
4. **Load Balancing**: Distributes computational load between Provider and Requestor

**Challenges:**

1. **Increased Complexity**: Requires more sophisticated algorithms and data exchange protocols  
2. **Privacy Concerns**: Necessitates sharing more detailed network information, which may be sensitive for some operators  
3. **Increased Computational Burden**: Both Provider and Requestor need to perform more complex calculations

**Use Cases:**

- Scenarios involving complex beam types (e.g., multipoint, coverage, or steerable beams)  
- When highly optimized service options are required, such as in military or emergency response situations  
- In trusted partnerships where detailed information sharing is acceptable

### 4.4.3. Hybrid and Adaptive Approaches

The Federation Architecture also supports hybrid and adaptive approaches that combine elements of both paths:

1. **Tiered Evaluation**: Initial service options are generated using the Provider-Only path, with the option to switch to the Provider \+ Requestor path for refined evaluation of promising candidates.  
     
2. **Dynamic Selection**: The approach is selected dynamically based on factors such as network complexity, privacy requirements, and computational resources available.  
     
3. **Partial Information Sharing**: Requestors can choose to share partial information, allowing for a middle ground between the two main approaches.

### 4.4.4. Considerations for Implementation

When implementing service option approaches, consider the following:

1. **Performance Optimization**: Implement efficient algorithms for service option generation and evaluation, especially for large-scale networks.  
     
2. **Privacy Safeguards**: Develop mechanisms to protect sensitive information when using the Provider \+ Requestor path, such as data anonymization or secure multi-party computation techniques.  
     
3. **Flexibility**: Design systems that can support both approaches, allowing operators to choose based on their specific requirements and constraints.  
     
4. **Temporal Considerations**: Ensure that both approaches can handle the temporal aspects of network resources, such as satellite visibility windows and time-varying service attributes.

By supporting multiple service option approaches, the Federation Architecture provides the flexibility to address a wide range of use cases, network configurations, and operator requirements. This adaptability is crucial for creating a robust and widely applicable Federation framework.

## 4.5. Temporal Considerations

The Federation Architecture places significant emphasis on temporal aspects, recognizing the dynamic nature of modern network environments, especially in space-based systems. Key temporal considerations include:

- **Time-Bound API Objects**: Most API objects (e.g., ServiceOptions, Services) incorporate time intervals to represent their validity periods. This allows for precise scheduling and management of resources that may only be available during specific windows.  
    
- **Dynamic Network Topology**: The architecture is designed to handle real-time changes in network topology. This is particularly crucial for satellite-based systems where connectivity options can change rapidly due to orbital dynamics.  
    
- **Temporal Service Attributes**: Service characteristics such as bandwidth, latency, and availability are represented as temporal values. This enables accurate modeling of service quality variations over time, accounting for factors like satellite position, atmospheric conditions, and network load.  
    
- **Granular Time Representation**: The architecture supports flexible time granularity, allowing for both fine-grained scheduling (e.g., second-by-second for LEO satellites) and coarser intervals for more stable network segments.  
    
- **Predictive Modeling**: By incorporating temporal data, the architecture facilitates predictive modeling of network performance and availability, enabling proactive resource allocation and service management.

This comprehensive approach to temporal considerations ensures that the Federation Architecture can accurately represent and manage the complex, time-varying nature of modern integrated space and terrestrial networks.

## 4.6. Security and Privacy Design Decisions

Security and privacy are paramount in the Federation Architecture, given the sensitive nature of network resource sharing and the potential involvement of multiple parties. Key security and privacy design decisions include:

- **Robust Authentication and Authorization**: The API incorporates built-in mechanisms for strong authentication and fine-grained authorization. This ensures that only verified entities can access resources and that their actions are strictly controlled based on their permissions.  
    
- **Flexible Information Sharing**: The architecture supports varying levels of information sharing based on trust relationships between partners. This allows operators to control the granularity and sensitivity of the data they expose to different Federation participants.  
    
- **Privacy-Preserving Techniques**: Especially in the peer-to-peer model, the architecture employs privacy-preserving techniques such as data minimization, anonymization, and secure multi-party computation. These methods allow for effective Federation while minimizing the exposure of sensitive network details.  
    
- **End-to-End Encryption Support**: For transmission of highly sensitive data, the architecture supports end-to-end encryption. This ensures that data remains confidential as it traverses multiple network segments and providers.  
    
- **Auditing and Logging**: Comprehensive auditing and logging capabilities are built into the architecture to support security monitoring, incident response, and compliance requirements.  
    
- **Segmentation and Isolation**: The architecture allows for logical segmentation of federated resources, ensuring that security breaches or issues in one part of the Federation do not compromise the entire system.

These security and privacy design decisions aim to create a trusted environment for Federation, balancing the need for secure, controlled access to network resources with the flexibility required for effective collaboration across diverse partners. The architecture provides a framework that can adapt to varying security requirements and regulatory landscapes across different regions and use cases.

## 4.7. Scalability and Performance Considerations

To ensure the architecture can handle large-scale Federations:

- Support for streaming updates is included to handle real-time changes in network conditions  
- Optimization techniques for service option generation and evaluation are implemented  
- The architecture allows for distributed processing in peer-to-peer scenarios to reduce central bottlenecks  
- Efficient data structures and algorithms for managing time-series data are employed

## 4.8. Interoperability Decisions

To facilitate seamless integration across diverse networks:

- Standardized API interfaces are defined for consistent interaction between different network operators  
- The architecture aligns with relevant industry standards (e.g., 3GPP, ETSI, ITU recommendations)  
- Support for common network protocols and data formats is included

These decisions aim to reduce integration barriers and enable a wide range of network operators to participate in the Federation.

## 4.9. Future Extensibility Considerations

The architecture is designed with future growth in mind:

- A modular design allows for the addition of new service types and network technologies  
- A versioning system is implemented to manage API evolution while maintaining backwards compatibility  
- Data models are designed to be extensible to accommodate emerging use cases and requirements

These considerations ensure that the Federation Architecture can adapt to new technologies, use cases, and operational models as the space and terrestrial network landscape continues to evolve.

# 5\. Open Issues and Future Work

This section outlines known limitations, areas requiring further development, and outstanding considerations for the Federation Architecture. These items represent active areas of discussion and development within the project.

## 5.1. API Design and Scope Considerations

### 5.1.1. API Complexity and Granularity

- Current approach of supporting all use cases in one general API may need reevaluation  
- Consider separation into per-layer or per-service type APIs:  
  - Dedicated Federation Service API for Optical sharing  
  - Dedicated Federation Service API for MHz sharing  
  - Separate API for networking requests  
  - Others as identified through implementation experience

### 5.1.2. Physical Layer Federation Support

- Current design focuses primarily on end-to-end network objectives  
- Need enhanced support for point-to-point Provider services (Optical, MHz)  
- Spectrum sharing scenarios require additional consideration:  
  - Provider compensation for spectrum non-use  
  - Dynamic spectrum allocation mechanisms  
  - Interference management protocols

## 5.2. Performance and Scaling Challenges

### 5.2.1. High-Frequency Operations

- Current API may not efficiently handle scenarios requiring rapid service updates  
- Need optimization for real-time tactical operations  
- Performance implications of frequent service scheduling updates

### 5.2.2. Stream Management Scalability

- Potential bottlenecks in continuous streaming of:  
  - Interconnection points  
  - Service status updates  
- Need for efficient client connection management  
- Resource optimization for high-update scenarios

## 5.3. Service Handoff and Continuity

### 5.3.1. Sequential Service Scheduling

- Lack of support for ServiceOption sequencing  
- Need for transactional handling of ServiceRequests  
- Time-series based accept/reject mechanisms

### 5.3.2. Make-Before-Break Support

- Current limitations in handoff strategy communication  
- Data continuity during asset transitions  
- Need for seamless service migration protocols

## 5.4. Business Operations Support

### 5.4.1. Service Discovery and Management

- Limited mechanisms for discovering generally available services  
- Need for long-term service availability forecasting  
- Regional service level support visibility

### 5.4.2. Commercial Operations

- Quoting and pricing mechanism requirements  
- Agreement formation and management  
- Usage tracking and billing support  
- Real-time pricing exposure  
- Resource consumption constraints and monitoring

## 5.5. Technical Considerations

### 5.5.1. Reliability and Redundancy

- Need for built-in redundancy mechanisms  
- High availability considerations  
- Regional synchronization of state  
- Fault tolerance protocols

### 5.5.2. SLA Management

- Enhanced SLA definition and tracking capabilities  
- Compliance monitoring and enforcement  
- Real-time SLA breach handling  
- Performance metric tracking

### 5.5.3. Corner Cases and Assumptions

- Off-nominal behavior documentation  
- Recovery processes for various failure scenarios  
- API temporal usage assumptions  
- Client session recovery protocols

## 5.6. Service Evolution and Support

### 5.6.1. Service Type Expansion

- Support for emerging service types  
- Specialized use case accommodation  
- Integration with new technologies:  
  - Delay Tolerant Networking  
  - IoT integration  
  - Edge computing services  
  - Dedicated tunneling options

### 5.6.2. Client Support Infrastructure

- Communication channel development  
- Issue reporting mechanisms  
- Feature request processes  
- Support ticket lifecycle management

## 5.7. Network Model Development

### 5.7.1. NMTS Enhancement

- Ongoing refinement based on real-world usage  
- Support for new use cases  
- Model expansion for emerging technologies  
- Integration with additional network types

### 5.7.2. Documentation and Guidelines

- Best practices development  
- Implementation guides  
- Use case examples  
- Integration patterns

These open issues and future work items represent active areas of development within the Federation Architecture. They are being addressed through ongoing research, development, and community feedback. Contributors are encouraged to participate in discussions and development efforts around these topics through the project's communication channels and development processes.

## 6.0. Contributors

The Outernet Council is thankful to the following individuals for their contributions to this reference architecture:

1. Brian Barritt  
2. Erik Kline  
3. David Mandle  
4. Paul Heninwolf  
5. Stefan Draskoci  
6. Helen Chou  
7. Scott Moeller  
8. Steve Nixon  
9. Michael Cheng

[^1]:  See the Licensed Shared Access (LSA) in Europe ([https://www.etsi.org/images/files/ETSIWhitePapers/ETSI-WinnForum-WPSpectrum\_sharing\_frameworks\_for\_temporary\_dynamic\_and\_flexible\_spectrum\_access\_for\_local\_private\_networks.pdf](https://www.etsi.org/images/files/ETSIWhitePapers/ETSI-WinnForum-WPSpectrum_sharing_frameworks_for_temporary_dynamic_and_flexible_spectrum_access_for_local_private_networks.pdf) ). The LSA required a standardized interface called the LSA1 interface. This interface supports the exchange of LSA Spectrum Resource Availability Information between the LSA Repository (LR) and LSA Controller (LC), as well as maintaining synchronization of this information. The standardized interface was necessary to ensure consistent communication between different components of the LSA system across various implementations. See also CBRS’s Spectrum Access System and Dynamic Spectrum Sharing (DSS) for LTE and 5G NR.

[image1]: <data:image/png;base64,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>