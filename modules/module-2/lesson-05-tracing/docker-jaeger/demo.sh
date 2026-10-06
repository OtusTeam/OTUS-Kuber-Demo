#!/bin/bash
set -e

source ../common.sh

echo "🚀 Starting Jaeger tracing demo..."

docker_compose down -v 2>/dev/null || true

docker_compose up -d

echo "⏳ Waiting for services to start..."
sleep 5

echo "✅ Services are ready!"
echo ""
echo "📊 Access points:"
echo "   Jaeger UI:  http://localhost:16686"
echo "   Demo App:   http://localhost:8080"
echo ""
echo "🧪 Generating test traces..."

for i in {1..10}; do
  curl -s -X POST http://localhost:8080/order >/dev/null && echo "Request $i sent"
  sleep 0.2
done

echo ""
echo "✅ Demo complete!"
echo "🔍 Open Jaeger UI: http://localhost:16686"
echo "   Service: demo-go-app"
echo ""
echo "To stop: docker_compose down"
