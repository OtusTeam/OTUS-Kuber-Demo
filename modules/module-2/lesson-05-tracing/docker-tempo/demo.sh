#!/bin/bash
set -e

echo "🚀 Starting Grafana Tempo tracing demo..."

docker-compose down -v 2>/dev/null || true

docker-compose up -d

echo "⏳ Waiting for services to start..."
sleep 8

echo "✅ Services are ready!"
echo ""
echo "📊 Access points:"
echo "   Grafana:    http://localhost:3000"
echo "   Tempo API:  http://localhost:3200"
echo "   Demo App:   http://localhost:8080"
echo ""
echo "🧪 Generating test traces..."

for i in {1..10}; do
  curl -s -X POST http://localhost:8080/order >/dev/null && echo "Request $i sent"
  sleep 0.2
done

echo ""
echo "✅ Demo complete!"
echo "🔍 Open Grafana: http://localhost:3000"
echo "   Navigate to: Explore → Tempo datasource"
echo "   Search for service: demo-go-app"
echo ""
echo "To stop: docker-compose down"
