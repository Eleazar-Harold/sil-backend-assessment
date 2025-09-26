#!/bin/bash

# Delete SIL Backend Assessment from Minikube
# This script removes all Kubernetes resources for the application

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
NAMESPACE="sil-backend-assessment"

echo -e "${BLUE}🗑️  Starting cleanup of SIL Backend Assessment from Minikube${NC}"

# Check if namespace exists
if ! kubectl get namespace ${NAMESPACE} > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Namespace ${NAMESPACE} does not exist. Nothing to clean up.${NC}"
    exit 0
fi

# Delete resources in reverse order to handle dependencies
echo -e "${YELLOW}📊 Deleting HPA...${NC}"
kubectl delete -f k8s/hpa.yaml --ignore-not-found=true

echo -e "${YELLOW}🌐 Deleting Ingress...${NC}"
kubectl delete -f k8s/ingress.yaml --ignore-not-found=true

echo -e "${YELLOW}🔗 Deleting Services...${NC}"
kubectl delete -f k8s/service.yaml --ignore-not-found=true

echo -e "${YELLOW}🚀 Deleting Deployment...${NC}"
kubectl delete -f k8s/deployment.yaml --ignore-not-found=true

echo -e "${YELLOW}🔴 Deleting Redis...${NC}"
kubectl delete -f k8s/redis.yaml --ignore-not-found=true

echo -e "${YELLOW}🗄️  Deleting PostgreSQL...${NC}"
kubectl delete -f k8s/postgres.yaml --ignore-not-found=true

echo -e "${YELLOW}🔐 Deleting Secrets...${NC}"
kubectl delete -f k8s/secret.yaml --ignore-not-found=true

echo -e "${YELLOW}📋 Deleting ConfigMap...${NC}"
kubectl delete -f k8s/configmap.yaml --ignore-not-found=true

# Wait for pods to be deleted
echo -e "${BLUE}⏳ Waiting for pods to be deleted...${NC}"
kubectl wait --for=delete pod -l app=sil-backend-assessment -n ${NAMESPACE} --timeout=120s || true
kubectl wait --for=delete pod -l app=postgres -n ${NAMESPACE} --timeout=120s || true
kubectl wait --for=delete pod -l app=redis -n ${NAMESPACE} --timeout=120s || true

# Delete namespace (this will delete any remaining resources)
echo -e "${YELLOW}📁 Deleting namespace...${NC}"
kubectl delete namespace ${NAMESPACE} --ignore-not-found=true

# Verify cleanup
echo -e "${BLUE}🔍 Verifying cleanup...${NC}"
if kubectl get namespace ${NAMESPACE} > /dev/null 2>&1; then
    echo -e "${RED}❌ Namespace still exists${NC}"
    exit 1
else
    echo -e "${GREEN}✅ Namespace deleted successfully${NC}"
fi

echo -e "${GREEN}🎉 Cleanup completed successfully!${NC}"
echo -e "${BLUE}📋 All resources have been removed from Minikube${NC}"
