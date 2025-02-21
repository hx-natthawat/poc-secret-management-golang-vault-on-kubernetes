# Zero to Hero: Vault Kubernetes Secret Management Setup Guide

This guide will walk you through setting up the Vault Kubernetes Secret Management PoC from scratch. Follow each step carefully to ensure a successful setup.

## Prerequisites Installation

### 1. Install Docker Desktop
```bash
# macOS (using Homebrew)
brew install --cask docker

# Start Docker Desktop from Applications
open -a Docker
```

### 2. Install Kubectl
```bash
# macOS (using Homebrew)
brew install kubectl

# Verify installation
kubectl version --client
```

### 3. Install Kind (Kubernetes in Docker)
```bash
# macOS (using Homebrew)
brew install kind

# Verify installation
kind version
```

### 4. Install Helm
```bash
# macOS (using Homebrew)
brew install helm

# Verify installation
helm version
```

### 5. Install Go
```bash
# macOS (using Homebrew)
brew install go

# Verify installation
go version  # Should show Go 1.22 or later
```

## Project Setup

### 1. Clone the Repository
```bash
# Clone the repository
git clone <repository-url>
cd poc-vault-go-kube
```

### 2. Create Kubernetes Cluster
```bash
# Create a new Kind cluster using our configuration
kind create cluster --config k8s/kind-config.yaml --name vault-demo

# Verify cluster is running
kubectl cluster-info
```

### 3. Create Namespaces
```bash
# Create required namespaces
kubectl create namespace vault
kubectl create namespace app

# Verify namespaces
kubectl get namespaces
```

## Vault Setup

### 1. Install Vault
```bash
# Add HashiCorp Helm repository
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update

# Install Vault
helm install vault hashicorp/vault \
  --namespace vault \
  -f k8s/vault/values.yaml

# Wait for Vault pod to be ready
kubectl -n vault wait --for=condition=ready pod -l app.kubernetes.io/name=vault
```

### 2. Initialize and Unseal Vault
```bash
# Initialize Vault
kubectl -n vault exec -it vault-0 -- vault operator init

# Save the output! You will need:
# - 5 Unseal Keys (save at least 3)
# - Initial Root Token

# Unseal Vault (run 3 times with different unseal keys)
kubectl -n vault exec -it vault-0 -- vault operator unseal
# Enter unseal key when prompted

# Verify Vault is unsealed
kubectl -n vault exec -it vault-0 -- vault status
```

### 3. Configure Vault
```bash
# Login to Vault
kubectl -n vault exec -it vault-0 -- vault login
# Enter the root token when prompted

# Enable Kubernetes authentication
kubectl -n vault exec -it vault-0 -- vault auth enable kubernetes

# Configure Kubernetes authentication
kubectl -n vault exec -it vault-0 -- vault write auth/kubernetes/config \
    kubernetes_host="https://kubernetes.default.svc.cluster.local:443"

# Enable KV v2 secret engines
kubectl -n vault exec -it vault-0 -- vault secrets enable -path=cluster-secrets kv-v2
kubectl -n vault exec -it vault-0 -- vault secrets enable -path=app-secrets kv-v2

# Create example secrets
kubectl -n vault exec -it vault-0 -- vault kv put cluster-secrets/database \
    url="postgresql://db.example.com:5432" \
    username="admin"

kubectl -n vault exec -it vault-0 -- vault kv put app-secrets/api-keys \
    sendgrid="sg_test_example" \
    stripe="sk_test_example"
```

## Application Deployment

### 1. Build and Load Docker Image
```bash
# Build the application
cd app
docker build -t vault-demo-app:latest .

# Load the image into Kind
kind load docker-image vault-demo-app:latest --name vault-demo
```

### 2. Deploy the Application
```bash
# Apply Kubernetes configurations
kubectl apply -f k8s/app/

# Wait for the application to be ready
kubectl -n app wait --for=condition=ready pod -l app=vault-demo-app

# Verify deployment
kubectl -n app get pods
```

## Verification

### 1. Check Vault UI Access
```bash
# Vault UI should be accessible at:
open http://localhost:30000

# Login using the root token saved earlier
```

### 2. Test Secret Retrieval
```bash
# Test cluster secrets endpoint
curl http://localhost:30001/cluster-secret/database

# Test application secrets endpoint
curl http://localhost:30001/app-secret/api-keys

# Test health endpoint
curl http://localhost:30001/health
```

## Common Issues and Solutions

### 1. Vault Pod Not Starting
```bash
# Check pod status
kubectl -n vault get pods
kubectl -n vault describe pod vault-0

# Check logs
kubectl -n vault logs vault-0
```

### 2. Application Authentication Issues
```bash
# Verify service account
kubectl -n app get serviceaccount

# Check pod logs
kubectl -n app logs -l app=vault-demo-app
```

### 3. Secret Access Issues
```bash
# Verify secret engine
kubectl -n vault exec -it vault-0 -- vault secrets list

# Check policy
kubectl -n vault exec -it vault-0 -- vault policy read app-policy
```

## Cleanup

### Remove Everything
```bash
# Delete the Kind cluster
kind delete cluster --name vault-demo

# Optional: Remove Docker images
docker rmi vault-demo-app:latest
```

## Next Steps

1. **Production Considerations**
   - Set up high availability
   - Configure proper backup strategies
   - Implement secret rotation
   - Set up monitoring and alerting

2. **Security Hardening**
   - Enable TLS
   - Implement stricter policies
   - Set up audit logging
   - Configure resource limits

3. **Advanced Features**
   - Dynamic secrets
   - Secret rotation
   - Transit encryption
   - Response wrapping

## Additional Resources

- [HashiCorp Vault Documentation](https://developer.hashicorp.com/vault/docs)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Go Documentation](https://golang.org/doc/)
- [Fiber Framework Documentation](https://docs.gofiber.io/)
