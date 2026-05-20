#!/bin/bash

# GoKube Deployment Script
set -e

echo "🚀 Starting GoKube deployment..."

# Check if minikube is running
if ! minikube status > /dev/null 2>&1; then
    echo "❌ Minikube is not running. Starting minikube..."
    minikube start
    minikube addons enable ingress
    echo "✅ Minikube started successfully"
fi

# Set minikube docker environment
echo "🔧 Setting up minikube docker environment..."
eval $(minikube docker-env)

# Build Docker image
echo "🏗️  Building Docker image..."
docker build -t gokube:latest .

# Create namespace if it doesn't exist
echo "📦 Creating namespace..."
kubectl create namespace gokube --dry-run=client -o yaml | kubectl apply -f -

# Deploy to Kubernetes
echo "🚀 Deploying to Kubernetes..."
kubectl apply -k k8s/

# Wait for deployment to be ready
echo "⏳ Waiting for deployment to be ready..."
kubectl wait --namespace gokube \
    --for=condition=ready pod \
    --selector=app=gokube \
    --timeout=120s

# Get service URL
echo "🌐 Getting service information..."
MINIKUBE_IP=$(minikube ip)
echo "Minikube IP: $MINIKUBE_IP"

# Show access methods
echo ""
echo "🎉 Deployment completed successfully!"
echo ""
echo "📋 Access your application:"
echo "1. NodePort Service: http://$MINIKUBE_IP:30080"
echo "2. Ingress (add to /etc/hosts):"
echo "   echo '$MINIKUBE_IP gokube.local' | sudo tee -a /etc/hosts"
echo "   Then visit: http://gokube.local"
echo ""
echo "🔍 Useful commands:"
echo "  kubectl get pods -n gokube"
echo "  kubectl get svc -n gokube"
echo "  kubectl get ingress -n gokube"
echo "  kubectl logs -f deployment/gokube -n gokube"
