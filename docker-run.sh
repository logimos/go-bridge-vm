#!/bin/bash

# Docker run script for myllm
set -e

echo "🐳 Building and running myllm Docker container..."

# Build the Docker image
echo "📦 Building Docker image..."
docker build -t myllm:latest .

# Run the container
echo "🚀 Starting myllm container..."
docker run -d \
  --name myllm \
  -p 8080:8080 \
  -e AI_PROVIDER=enhanced_local \
  -e INTENT_CONFIG_PATH=configs/personal_assistant.json \
  -e AI_TEMPERATURE=0.1 \
  -e AI_MAX_TOKENS=1000 \
  --restart unless-stopped \
  myllm:latest

echo "✅ Container started successfully!"
echo "🌐 Server is running on http://localhost:8080"
echo ""
echo "📋 Useful commands:"
echo "  - View logs: docker logs myllm"
echo "  - Stop container: docker stop myllm"
echo "  - Remove container: docker rm myllm"
echo "  - Test API: curl http://localhost:8080/api/v1/debug"
echo ""
echo "🔧 To use different environment variables, edit this script or use:"
echo "   docker run -e AI_PROVIDER=openai -e OPENAI_API_KEY=your-key ..." 