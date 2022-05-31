# opentelemetry go example

OpenTelemetry is a set of APIs, SDKs, tooling and integrations that are designed for the creation and management of 
telemetry data such as traces, metrics, and logs.

![Reference Architecture](/docs/pics/otel-reference-architecture.svg)

## Which problem it solves

Distributed architectures introduce a variety of operational challenges including how to solve availability and 
performance issues quickly.

Telemetry data is needed to power observability products. Traditionally, telemetry data has been provided by either 
open-source projects or commercial vendors. With a lack of standardization, the net result is the lack of data 
portability and the burden on the user to maintain the instrumentation.

OpenTelemetry is not an observability back-end like Jaeger or Prometheus.

(source: https://opentelemetry.io/docs/concepts/what-is-opentelemetry/)

## Example

A very simplified example.

![example](/docs/pics/example-services.png)

### Overview


| Service         |                                                                                |
|-----------------|--------------------------------------------------------------------------------|
| consumer        | Interact with the user management service. It is doing random domain events.   |
| user-management | This service manage the state of user data                                     |

### Domain events

| Events      |                |
|-------------|----------------|
| user.list   | List user      |
| user.create | Creates a user |
| user.delete | Delete a user  |

Events are not single requests. We should consider that an event have more impact. For example the user create maybe 
is interesting for an analytics team. Or the delete event will have an impact on GDPR.

### Further info:
http://www.inanzzz.com/index.php/post/4qes/implementing-opentelemetry-and-jaeger-tracing-in-golang-http-api
