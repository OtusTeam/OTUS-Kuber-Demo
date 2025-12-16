# Grafana Tempo Tracing Demo

Example demonstrating distributed tracing with Grafana Tempo and visualization in Grafana.

## Components

- **Grafana Tempo**: Scalable distributed tracing backend
- **Grafana**: Visualization and querying with TraceQL
- **OpenTelemetry Collector**: Receives traces from app and exports to Tempo
- **Demo App**: Go application with OpenTelemetry instrumentation

## Quick Start

```bash
./demo.sh
```

## Manual Setup

```bash
# Start services
docker-compose up -d

# Wait for services
sleep 8

# Generate traces
for i in {1..10}; do
  curl -X POST http://localhost:8080/order
  sleep 0.2
done

# Open Grafana
open http://localhost:3000
```

## Access Points

- **Grafana**: http://localhost:3000 (no login required)
- **Tempo API**: http://localhost:3200
- **Demo App**: http://localhost:8080

## View Traces in Grafana

1. Open http://localhost:3000
2. Navigate to: Explore (compass icon in sidebar)
3. Select datasource: "Tempo"
4. Search methods:
   - **Search**: Service name = `demo-go-app`
   - **TraceQL**: `{ service.name = "demo-go-app" }`
5. Click on a trace to see flame graph and span details

## TraceQL Examples

```traceql
# Find all traces for demo-go-app
{ service.name = "demo-go-app" }

# Find traces with errors
{ service.name = "demo-go-app" && status = error }

# Find slow traces (>100ms)
{ service.name = "demo-go-app" } | duration > 100ms
```

## Architecture

```
┌──────────┐   OTLP/HTTP   ┌──────────────┐   OTLP/gRPC   ┌─────────┐
│ Demo App │──────────────▶│ OTel         │──────────────▶│ Tempo   │
│          │               │ Collector    │               │         │
└──────────┘               └──────────────┘               └─────────┘
                                                                │
                                                                ▼
                                                           ┌─────────┐
                                                           │ Grafana │
                                                           │   UI    │
                                                           └─────────┘
                                                        localhost:3000
```

## Features

- **TraceQL**: Powerful query language for traces
- **Service Graph**: Visualize service dependencies
- **Metrics from Traces**: Generate metrics from span data
- **Scalable Storage**: Optimized for high-volume tracing

## Cleanup

```bash
docker-compose down -v
```
