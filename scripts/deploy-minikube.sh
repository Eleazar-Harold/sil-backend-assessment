#!/bin/bash

# Deploy SIL Backend Assessment to Minikube
# This script builds the Docker image, loads it into Minikube, and deploys the application

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
IMAGE_NAME="sil-backend-assessment"
IMAGE_TAG="latest"
NAMESPACE="sil-backend-assessment"

echo -e "${BLUE}🚀 Starting Minikube deployment for SIL Backend Assessment${NC}"

# Check if Minikube is running
if ! minikube status > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Minikube is not running. Starting Minikube...${NC}"
    minikube start --memory=4096 --cpus=2 --disk-size=20g
    echo -e "${GREEN}✅ Minikube started successfully${NC}"
fi

# Set Docker environment to use Minikube's Docker daemon
echo -e "${BLUE}🔧 Setting up Docker environment for Minikube...${NC}"
eval $(minikube docker-env)

# Build the Docker image
echo -e "${BLUE}🏗️  Building Docker image...${NC}"
docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .
echo -e "${GREEN}✅ Docker image built successfully${NC}"

# Load image into Minikube
echo -e "${BLUE}📦 Loading image into Minikube...${NC}"
minikube image load ${IMAGE_NAME}:${IMAGE_TAG}
echo -e "${GREEN}✅ Image loaded into Minikube${NC}"

# Create namespace if it doesn't exist
echo -e "${BLUE}📁 Creating namespace...${NC}"
kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
echo -e "${GREEN}✅ Namespace created/verified${NC}"

# Apply Kubernetes manifests
echo -e "${BLUE}🚀 Deploying application to Kubernetes...${NC}"

# Apply in order to handle dependencies
echo -e "${YELLOW}📋 Applying ConfigMap and Secrets...${NC}"
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml

echo -e "${YELLOW}🗄️  Deploying PostgreSQL...${NC}"
kubectl apply -f k8s/postgres.yaml

echo -e "${YELLOW}🔴 Deploying Redis...${NC}"
kubectl apply -f k8s/redis.yaml

# Wait for databases to be ready
echo -e "${BLUE}⏳ Waiting for databases to be ready...${NC}"
kubectl wait --for=condition=ready pod -l app=postgres -n ${NAMESPACE} --timeout=300s
kubectl wait --for=condition=ready pod -l app=redis -n ${NAMESPACE} --timeout=300s
echo -e "${GREEN}✅ Databases are ready${NC}"

echo -e "${YELLOW}🚀 Deploying main application...${NC}"
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

echo -e "${YELLOW}🌐 Setting up ingress...${NC}"
kubectl apply -f k8s/ingress.yaml

echo -e "${YELLOW}📊 Setting up HPA...${NC}"
kubectl apply -f k8s/hpa.yaml

# Wait for application to be ready
echo -e "${BLUE}⏳ Waiting for application to be ready...${NC}"
kubectl wait --for=condition=ready pod -l app=sil-backend-assessment -n ${NAMESPACE} --timeout=300s
echo -e "${GREEN}✅ Application is ready${NC}"

# Get service URLs
echo -e "${BLUE}🌍 Getting service information...${NC}"
MINIKUBE_IP=$(minikube ip)
NODEPORT=$(kubectl get service sil-backend-nodeport -n ${NAMESPACE} -o jsonpath='{.spec.ports[0].nodePort}')

echo -e "${GREEN}🎉 Deployment completed successfully!${NC}"
echo -e "${BLUE}📋 Service Information:${NC}"
echo -e "   🌐 NodePort URL: http://${MINIKUBE_IP}:${NODEPORT}"
echo -e "   🔍 Health Check: http://${MINIKUBE_IP}:${NODEPORT}/api/health"
echo -e "   📚 API Docs: http://${MINIKUBE_IP}:${NODEPORT}/docs"
echo -e "   🎮 GraphQL Playground: http://${MINIKUBE_IP}:${NODEPORT}/graphql/playground"

# Show pods status
echo -e "${BLUE}📊 Pod Status:${NC}"
kubectl get pods -n ${NAMESPACE}

echo -e "${BLUE}🔗 Services:${NC}"
kubectl get services -n ${NAMESPACE}

echo -e "${BLUE}📈 Ingress:${NC}"
kubectl get ingress -n ${NAMESPACE}

# Optional: Open browser to health check
if command -v open > /dev/null 2>&1; then
    echo -e "${YELLOW}🌐 Opening browser to health check...${NC}"
    open "http://${MINIKUBE_IP}:${NODEPORT}/api/health"
elif command -v xdg-open > /dev/null 2>&1; then
    echo -e "${YELLOW}🌐 Opening browser to health check...${NC}"
    xdg-open "http://${MINIKUBE_IP}:${NODEPORT}/api/health"
fi

echo -e "${GREEN}✅ Minikube deployment completed!${NC}"
echo -e "${YELLOW}💡 To view logs: kubectl logs -f deployment/sil-backend-assessment -n ${NAMESPACE}${NC}"
echo -e "${YELLOW}💡 To delete deployment: ./scripts/delete-minikube.sh${NC}"
