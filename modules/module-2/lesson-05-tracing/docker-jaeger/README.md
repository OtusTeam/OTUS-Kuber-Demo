# Jaeger Tracing Demo

Simple example demonstrating distributed tracing with Jaeger.

## Components

- **Jaeger**: All-in-one tracing backend (UI + collector)
- **OpenTelemetry Collector**: Receives traces from app and exports to Jaeger
- **Demo App**: Go application with OpenTelemetry instrumentation

## Quick Start

```bash
./demo.sh
```

This will:
1. Start all services
2. Generate test traces
3. Open Jaeger UI

## Manual Setup

```bash
# Start services
docker-compose up -d

# Wait for services
sleep 5

# Generate traces
for i in {1..10}; do
  curl -X POST http://localhost:8080/order
  sleep 0.2
done

# Open Jaeger UI
open http://localhost:16686
```

## Access Points

- **Jaeger UI**: http://localhost:16686
- **Demo App**: http://localhost:8080

## View Traces

1. Open http://localhost:16686
2. Select service: `demo-go-app`
3. Click "Find Traces"
4. Click on a trace to see spans

## Architecture

```
┌──────────┐   OTLP/HTTP   ┌──────────────┐   OTLP/gRPC   ┌─────────┐
│ Demo App │──────────────▶│ OTel         │──────────────▶│ Jaeger  │
│          │               │ Collector    │               │         │
└──────────┘               └──────────────┘               └─────────┘
                                                                │
                                                                ▼
                                                           Jaeger UI
                                                        localhost:16686
```

## Cleanup

```bash
docker-compose down -v
```
