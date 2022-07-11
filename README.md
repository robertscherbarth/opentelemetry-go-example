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


| Service         |                                                                              |
|-----------------|------------------------------------------------------------------------------|
| consumer        | Interact with the user management service. It is doing random domain events. |
| user-management | This service manage the state of user data                                   |
| analytics       | Service that gets data from a user.create event and enrich span              |

### Domain events

| Events      |                |
|-------------|----------------|
| user.list   | List user      |
| user.create | Creates a user |
| user.delete | Delete a user  |

Events are not single requests. It should be considered that an event have more impact. For example a `user.create` event 
maybe is interesting for an analytics team. Or the `user.delete` event will have an impact on GDPR.

## Setup

### Jaeger Setup

This setup runs out of the box. It includes a local jaeger that is used to visualize the traces. 

#### Run the setup

```
make jaeger-build && make jaeger-up  
```

### Honeycomb Setup

This docker-compose setup will export the data to https://www.honeycomb.io/.

To enable this you need to add a `.env-honeycomb` file in this folder this need an environment definition with this env key `HONEYCOMB_API_KEY`

example:
```
HONEYCOMB_API_KEY=1234567890
```

#### Run the setup

```
make honeycomb-build && make honeycomb-up
```


### Further info:
http://www.inanzzz.com/index.php/post/4qes/implementing-opentelemetry-and-jaeger-tracing-in-golang-http-api
