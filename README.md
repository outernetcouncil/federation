# Federation API: Unifying Space and Terrestrial Networks

**Seamlessly integrate diverse network resources for unparalleled global connectivity**

## Key Benefits of the Federation API

1. **Enhanced Operational Flexibility**: Dynamically allocate resources across federated networks to adapt to changing needs and ensure continuous connectivity in any situation.

2. **Expanded Market Reach**: Unlock new opportunities by making your network resources available to a wider range of customers and partners, including government, military, and commercial users.

3. **Optimized Resource Utilization**: Maximize the value of your network investments through intelligent resource sharing and dynamic capacity allocation, improving efficiency and reducing costs.

4. **Accelerated Innovation**: Rapidly deploy new services and enter markets by leveraging existing infrastructure and partnerships, reducing time-to-market for new offerings.

5. **Improved Global Connectivity**: Ensure reliable communications anywhere, anytime through seamless integration of terrestrial, aerial, and space-based networks, enhancing disaster response and emergency communications capabilities.

## Key Terms

- **Requestor**: An entity seeking network services or resources from federated partners.
- **Provider**: An entity offering network services or resources to federated partners.
- **Interconnection Point**: A physical or logical point where two federated networks can connect and exchange traffic.
- **Service Option**: A potential service offering from a provider, including details on availability, performance, and cost.
- **Availability**: The time periods during which a network resource or service is accessible for use.
- **Reachability**: The set of destinations or network prefixes that can be accessed through a given network resource or service.
- **Federation**: The act of combining multiple independent networks to create a larger, more capable network ecosystem.

## How the Federation API Works

The Federation API enables seamless collaboration between diverse network operators by providing a standardized interface for resource advertising, discovery, and provisioning. Providers publish their available resources and interconnection points, which requestors can then discover and evaluate. Requestors can submit service requests, specifying their requirements and preferences. The API facilitates negotiation between parties, allowing for dynamic service provisioning and resource allocation across federated networks. Throughout the service lifecycle, the API enables real-time status updates and modifications, ensuring optimal performance and efficiency.

## Federation API Paradigms

### Operational Security Paradigm
In this approach, requestors prioritize security by minimizing exposure of their network details. They first explore the provider's network using `StreamInterconnectionPoints`, then make a `ScheduleServiceRequest` with only essential information. This paradigm is ideal for military or sensitive commercial operations, allowing requestors to maintain control over their network information while still leveraging federated resources.

### Provider Solution Offloading Paradigm
This paradigm allows requestors to expose one or more of their assets to the provider, enabling the provider to determine the optimal connectivity solution. By using a more comprehensive `ScheduleServiceRequest`, requestors can leverage the provider's expertise and network knowledge to achieve the best possible service configuration. This approach is suitable for commercial partnerships where trust is established and optimizing performance is the primary goal.

## Real-World Use Cases

1. **Military Operations**: A defense agency uses the Federation API to securely access commercial satellite capacity during a mission, enhancing communication capabilities without compromising operational security. The API enables rapid service provisioning and seamless integration with existing military networks.

2. **GEO/LEO Hybrid Services**: A GEO satellite operator partners with a LEO constellation provider to offer low-latency, high-capacity services globally. The Federation API facilitates dynamic resource allocation between the two networks, optimizing performance and coverage for end-users.

3. **Cellular Coverage Expansion**: A major cellular provider leverages the Federation API to integrate satellite connectivity into its network, extending coverage to rural and remote areas. The API enables seamless handover between terrestrial and satellite resources, ensuring consistent user experience.

4. **HAPS Network Integration**: A HAPS operator uses the Federation API to integrate its high-altitude platforms with existing satellite and terrestrial networks. This enables the HAPS provider to offer unique mid-altitude connectivity options and expand its service portfolio rapidly.

5. **Global IoT Deployment**: An IoT service provider utilizes the Federation API to create a global network leveraging multiple satellite constellations and terrestrial networks. The API enables efficient device management and data routing across diverse network segments.

## Example Call Flows

### Operational Security Paradigm

1. Requestor calls `StreamInterconnectionPoints` to explore Provider's network
2. Requestor processes received InterconnectionPoints and determines suitable options
3. Requestor calls `ScheduleService` with minimal network details and preferred interconnection points
4. Provider processes request and returns a `ScheduleServiceResponse` with a unique `service_id`
5. Requestor calls `MonitorServices` to receive status updates for the service
6. Provider sends `ServiceStatus` updates through the `MonitorServices` stream
7. When service is no longer needed, Requestor calls `CancelService` with the `service_id`
8. Provider confirms cancellation and terminates the service

### Provider Solution Offloading Paradigm

1. Requestor calls `ListServiceOptions` with details of their network assets and requirements
2. Provider returns a stream of `ServiceOption` messages
3. Requestor evaluates received options and selects preferred solution
4. Requestor calls `ScheduleService` with chosen `ServiceOption` and additional details
5. Provider processes request and returns a `ScheduleServiceResponse` with a unique `service_id`
6. Requestor calls `MonitorServices` to receive status updates for the service
7. Provider sends `ServiceStatus` updates through the `MonitorServices` stream
8. If changes are needed, Requestor can update the service using `ScheduleService` with the existing `service_id`
9. When service is no longer needed, Requestor calls `CancelService` with the `service_id`
10. Provider confirms cancellation and terminates the service

## Detailed API Guide

[Detailed API Guide](APIGUIDE.md)

## Contributing

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE.txt) file for details.