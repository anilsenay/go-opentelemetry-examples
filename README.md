# go-opentelemetry-examples

I implemented some example cases for [OpenTelemetry](https://opentelemetry.io/) by using Golang.

Go Packages Used:
- [opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go): Go implementation of OpenTelemetry. It provides a set of APIs to directly measure performance and behavior of your software and send this data to observability platforms.
- [fiber](https://github.com/gofiber/fiber): Express inspired web framework
- [gorm](https://github.com/go-gorm/gorm): ORM library

There is 5 different example cases in seperated branches:
- [OpenTelemetry + Jaeger](#opentelemetry--jaeger): Basic example by using OpenTelemetry + [Jaeger](https://www.jaegertracing.io/)
- [OpenTelemetry + Jaeger + Distributed Services](#opentelemetry--jaeger--distributed-services): A distributed services example by using OpenTelemetry + [Jaeger](https://www.jaegertracing.io/)
- [OpenTelemetry + Jaeger + Elasticsearch](#opentelemetry--jaeger--elasticsearch): Using Elasticsearch as Jaeger span storage
- [OpenTelemetry + Elasticsearch APM](#opentelemetry--elasticsearch-apm): Basic example of using [Elasticsearch APM](https://www.elastic.co/observability/application-performance-monitoring) instead of Jaeger
- [(METRICS) OpenTelemetry + Jaeger + Prometheus Go Client](): Basic metric example by using OpenTelemetry + Jaeger + Prometheus Go Client
- [OpenTelemetry + Jaeger + Distributed Services + AWS Trace ID (X-Amzn-Trace-Id)](#opentelemetry--jaeger--distributed-services--aws-trace-id-x-amzn-trace-id): Example of passing the received `X-Amzn-Trace-Id` to routes

## OpenTelemetry + Jaeger

**Branch:** [master](https://github.com/anilsenay/go-opentelemetry-examples/tree/master): OpenTelemetry + [Jaeger](https://www.jaegertracing.io/)

_NOTE_: Jaeger uses inmemory storage in this case. So its not recommended for production. In production you should use a persistence storage. [Check here for more details](https://www.jaegertracing.io/docs/1.49/deployment/#span-storage-backends).

In this example tracing route is: `Handler(Fiber)` -> `Service` -> `Repository(Gorm)` -> `Database`
Fiber and Gorm middlewares do the trick for you, but you can add a new custom trace span just like I did in `todo_service.go`. **Its important to propagate context to sub routes.**

![Jaeger UI](https://github.com/anilsenay/go-opentelemetry-examples/assets/1047345/190d4eb2-6055-489c-b004-6febbd91f069)

## OpenTelemetry + Jaeger + Distributed Services

**Branch:** [distributed](https://github.com/anilsenay/go-opentelemetry-examples/tree/distributed): OpenTelemetry + [Jaeger](https://www.jaegertracing.io/) + Distributed Services

_NOTE_: Jaeger uses inmemory storage in this case. So its not recommended for production. In production you should use a persistence storage. [Check here for more details](https://www.jaegertracing.io/docs/1.49/deployment/#span-storage-backends).
In this example tracing route is: `TodoService1` -> `TodoService2` -> `TodoService3` -> `Database`

As you see in image, we can see the whole trace among services while sending request to another service by passing the context

![Jaeger UI](https://github.com/anilsenay/go-opentelemetry-examples/assets/1047345/26bca295-791a-45c7-925d-03218cd5e919)

## OpenTelemetry + Jaeger + Elasticsearch

**Branch:** [jaeger-elastic-storage](https://github.com/anilsenay/go-opentelemetry-examples/tree/jaeger-elastic-storage): OpenTelemetry + [Jaeger](https://www.jaegertracing.io/) + [Elasticsearch](https://www.elastic.co/)

Jaeger uses Elasticsearch as span storage in this case. The only difference is setting Jaeger configuration in docker-compose.yaml file.
For this example, I used Elastic Cloud instead of local Elastic instance.

**IMPORTANT**: When I tried, Jaeger failed with Elastic version 8. So I used Elastic version 7. If you have any trouble, you may try using version 7

## OpenTelemetry + Elasticsearch APM

**Branch:** [elastic](https://github.com/anilsenay/go-opentelemetry-examples/tree/elastic): OpenTelemetry + [Elasticsearch APM](https://www.elastic.co/observability/application-performance-monitoring)

![Elasticsearch APM Traces](https://github.com/anilsenay/go-opentelemetry-examples/assets/1047345/006abab3-0c48-462e-b9c8-55678a833270)

## (METRICS) OpenTelemetry + Jaeger + Prometheus Go Client

**Branch:** [metrics](https://github.com/anilsenay/go-opentelemetry-examples/tree/metrics): OpenTelemetry + [Jaeger](https://www.jaegertracing.io/) + [Prometheus Go Client](https://github.com/prometheus/client_golang)

_NOTE:_ This example does not include Prometheus instance for collecting metrics. It only serve metrics from /metrics endpoint in application.

![Prometheus Metrics Endpoint](https://github.com/anilsenay/go-opentelemetry-examples/assets/1047345/bd669285-9981-440f-b97b-930412ea61ec)

## OpenTelemetry + Jaeger + Distributed Services + AWS Trace ID (X-Amzn-Trace-Id)

**Branch:** [distributed-amazon-trace-id](https://github.com/anilsenay/go-opentelemetry-examples/tree/distributed-amazon-trace-id): OpenTelemetry + [Jaeger](https://www.jaegertracing.io/) + Distributed Services + [X-Amzn-Trace-Id](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-request-tracing.html)

I added this example for the case if you using `X-Amzn-Trace-Id` in somewhere in your system, you can easily filter the route of request with this tag.

For exampe if service received a `X-Amzn-Trace-Id` as `Root=1-67891233-abcdef012345678912345678` from AWS Load Balancer:
![Jaeger UI](https://github.com/anilsenay/go-opentelemetry-examples/assets/1047345/b5fa6d36-9ddf-48a3-a806-187ed180371e)
