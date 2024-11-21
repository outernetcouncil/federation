<img src="fedlogo2.png" width="681" height="257" alt="Description">
# Federation Reference Architecture v0.1

# Table of Contents

<details>
<summary>Table of Contents</summary>
  
1. [Introduction](#1-introduction)  
1.1 [Purpose of the Federation Architecture Specification](#11-purpose-of-the-federation-architecture-specification)  
1.2 [Scope of the Document](#12-scope-of-the-document)  
1.3 [Background and Context](#13-background-and-context)  
1.4 [Key Use Cases](#14-key-use-cases)  
1.5 [Intended Audience](#15-intended-audience)  
1.6 [Related Documents and Resources](#16-related-documents-and-resources)  
1.7 [Versioning and Updates](#17-versioning-and-updates)  
2. [Terms and Definitions](#2-terms-and-definitions)  
3. [Architecture Description Overview](#3-architecture-description-overview)  
3.1 [Federation Architecture Vision](#31-federation-architecture-vision)  
3.2 [Key Principles and Concepts](#32-key-principles-and-concepts)  
3.3 [Federation Models](#33-federation-models)  
3.4 [API Structure and Core Resources](#34-api-structure-and-core-resources)  
3.5 [Federation Workflow Overview](#35-federation-workflow-overview)  
3.6 [Temporal Considerations](#36-temporal-considerations)  
3.7 [Security and Privacy Considerations](#37-security-and-privacy-considerations)  
3.8 [Scalability and Performance](#38-scalability-and-performance)  
3.9 [Interoperability and Standards](#39-interoperability-and-standards)  
3.10 [Future Extensibility](#310-future-extensibility)  
4. [Architecture Rationale](#4-architecture-rationale)  
4.1 [Key Architectural Decisions](#41-key-architectural-decisions)  
4.2 [Trade-offs and Alternatives](#42-trade-offs-and-alternatives)  
4.3 [Federation Models](#43-federation-models)  
4.4 [Service Option Approaches](#44-service-option-approaches)  
4.5 [Temporal Considerations](#45-temporal-considerations)  
4.6 [Security and Privacy Design Decisions](#46-security-and-privacy-design-decisions)  
4.7 [Scalability and Performance Considerations](#47-scalability-and-performance-considerations)  
4.8 [Interoperability Decisions](#48-interoperability-decisions)  
4.9 [Future Extensibility Considerations](#49-future-extensibility-considerations)
5. [Open Issues and Future Work](#5-open-issues-and-future-work)
6. [Contributors](#6-contributors)
</details>

# 1\. Introduction

## 1.1 Purpose of the Federation Architecture Specification

Federation is a protocol that connects networks together in a way that unlocks new capabilities and enables seamless interoperability between participant networks. Bringing together various interoperable networks with their unique strengths and capabilities enables more complex and capable user experiences. As an open community led by its contributors, Federation’s architecture is evolving as network technologies are brought in to adapt to evolving customer needs. The Outernet Council’s development and evolution of Federation is guided by how the network can serve the user’s needs along the following themes:

* Reducing costs  
* Improving choice  
* Network resilience  
* Network security  
* Enabling new capabilities

Federation is part of a larger vision of a connected, secure, resilient, network of communications satellites provisioned by a diverse ecosystem of operators across multiple countries – the Outernet. The Outernet has the potential to unify space and terrestrial network architectures across diverse network segments, including land, sea, air and space.

The Outernet envisions broad interoperability between commercial and government satellite constellations, ground stations and user terminals.

<img src="outernetref.png" width="690.2" height="452.2" alt="Description">
</div>

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

In the most basic form, Federation enables cooperation between two disparate networks to overcome the outage/unavailability of a node -- allowing two networks to fulfil a service request that wouldn't otherwise be possible.

<img src="ref-ex1.png" width="760.2" height="433.3" alt="Description">
</div>

When leveraged more extensively, Federation has the potential to dramatically extend the capabilities and resillience of multiple archetypes of next generation networks:

### Space-Relay Interconnect and Transit

* Government satellite interconnecting in space (via RF or optical) with a commercial constellation for on-demand network transit to a given destination

### Ground Segment as a Service

* Dynamic reservation of third party ground stations (optical or RF) for on-demand network transit to specified points of presence, ops centers, or storage

### “Agile” (Multi-Constellation) User Terminals

* Multiple projects are underway to develop "agile" user terminals with multiple physical or virtualized modems that share a single aperture and are capable of dynamically switching between multiple commercial or military SATCOM systems  
* Federation enables just-in time on-demand ordering and provisioning of beams and transit for fixed terrestrial, land mobile, airborne, or maritime terminals to roam between these networks

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
  - [APIGUIDE.md](http://APIGUIDE.md](https://github.com/outernetcouncil/federation/blob/main/APIGUIDE.md)), a detailed guide to supporting the Federation API  
  - [README.md](http://README.md](https://github.com/outernetcouncil/federation/blob/main/README.md)), a succinct motivation summary of the Federation API  
  - [Federation.proto]([http://federation.proto](https://github.com/outernetcouncil/federation/blob/main/outernet/federation/v1alpha/federation.proto)) Protocol Buffer definitions for the Federation Service  
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

<img src="requestor.png" width="338" height="583" alt="Description">
</div>

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
4. Nihar Agrawal
5. Paul Heninwolf  
6. Stefan Draskoci  
7. Helen Chou  
8. Scott Moeller  
9. Steve Nixon  
10. Michael Cheng

[^1]:  See the Licensed Shared Access (LSA) in Europe ([https://www.etsi.org/images/files/ETSIWhitePapers/ETSI-WinnForum-WPSpectrum\_sharing\_frameworks\_for\_temporary\_dynamic\_and\_flexible\_spectrum\_access\_for\_local\_private\_networks.pdf](https://www.etsi.org/images/files/ETSIWhitePapers/ETSI-WinnForum-WPSpectrum_sharing_frameworks_for_temporary_dynamic_and_flexible_spectrum_access_for_local_private_networks.pdf) ). The LSA required a standardized interface called the LSA1 interface. This interface supports the exchange of LSA Spectrum Resource Availability Information between the LSA Repository (LR) and LSA Controller (LC), as well as maintaining synchronization of this information. The standardized interface was necessary to ensure consistent communication between different components of the LSA system across various implementations. See also CBRS’s Spectrum Access System and Dynamic Spectrum Sharing (DSS) for LTE and 5G NR.
