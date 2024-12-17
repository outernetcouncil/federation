# Federation API: Unifying Space and Terrestrial Networks

**Seamlessly integrate diverse network resources for unparalleled global connectivity**

[Glossary of Terms](docs/GLOSSARY.md)

## Purpose of the Federation API

[Federation](docs/GLOSSARY.md#key-concepts) is a protocol that connects networks together in a way that unlocks new capabilities and enables seamless interoperability between participant networks. Bringing together various interoperable networks with their unique strengths and capabilities enables more complex and capable user experiences. As an open community led by its contributors, Federation’s architecture is evolving as network technologies are brought in to adapt to evolving customer needs. The Outernet Council’s development and evolution of Federation is guided by how the network can serve the user’s needs along the following themes:

* Reducing costs
* Improving choice
* Network resilience
* Network security
* Enabling new capabilities

Federation is part of a larger vision of a connected, secure, resilient, network of communications satellites provisioned by a diverse ecosystem of operators across multiple countries – the Outernet. The Outernet has the potential to unify space and terrestrial network architectures across diverse network segments, including land, sea, air and space.

The Outernet envisions broad interoperability between commercial and government satellite constellations, ground stations and user terminals.

<img src="docs/outernetref.png" width="690.2" height="452.2" alt="Description">
</div>

The primary purposes of this specification are to:

1. Establish a standardized approach for integrating diverse network resources
2. Enable dynamic allocation of resources across federated networks
3. Facilitate seamless communication between terrestrial, aerial, and space-based systems
4. Provide a framework for expanding market reach and optimizing resource utilization
5. Accelerate innovation in global connectivity solutions

Additionally, this specification directly addresses the complexities of spatio-temporal asset utilization, such as scheduling assets with complex time-dependent availability and capacity.

## Background and Context

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

## Key Use Cases

In the most basic form, Federation enables cooperation between two disparate networks to overcome the outage/unavailability of a node -- allowing two networks to fulfil a service request that wouldn't otherwise be possible.

<img src="docs/ref-ex1.png" width="760.2" height="433.3" alt="Description">
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

## Detailed API Guide & Reference Architecture

* [Detailed API Guide](APIGUIDE.md)
* [Reference Architecture](https://github.com/outernetcouncil/federation/blob/main/docs/REFERENCE.md)

## Contributing

This project requires [Bazel](https://bazel.build/) for building and testing.
Please refer to the [Bazel documentation](https://docs.bazel.build/versions/main/install.html) for installation instructions.

To build all targets of the project, run the following command:

```bash
bazel build //...
```

To just build the protobuf definitions and the go definitions, run:

```bash
bazel build federation_proto
bazel build federation_go_proto
```

You can also see how the protobufs can be used in a go bazel project by looking at the examples folder.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE.txt) file for details.
