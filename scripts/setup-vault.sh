#!/bin/bash

# Exit on error
set -e

echo "🚀 Setting up Vault demo environment..."

# Check prerequisites
command -v docker >/dev/null 2>&1 || { echo "❌ Docker is required but not installed. Aborting." >&2; exit 1; }
command -v kind >/dev/null 2>&1 || { echo "❌ Kind is required but not installed. Aborting." >&2; exit 1; }
command -v kubectl >/dev/null 2>&1 || { echo "❌ kubectl is required but not installed. Aborting." >&2; exit 1; }
command -v helm >/dev/null 2>&1 || { echo "❌ Helm is required but not installed. Aborting." >&2; exit 1; }

echo "✅ Prerequisites checked"

# Create Kind cluster
echo "📦 Creating Kind cluster..."
kind create cluster --config kind/kind-config.yaml --name vault-demo

# Create namespaces
echo "🏗️ Creating namespaces..."
kubectl create namespace vault
kubectl create namespace app

# Add Helm repo
echo "📚 Adding HashiCorp Helm repo..."
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update

# Install Vault
echo "🔒 Installing Vault..."
helm install vault hashicorp/vault \
  --namespace vault \
  -f k8s/vault/values.yaml

# Wait for Vault pod
echo "⏳ Waiting for Vault pod to be ready..."
kubectl wait --for=condition=Ready pod/vault-0 -n vault --timeout=120s

echo "
✨ Setup complete! Next steps:

1. Initialize Vault:
   kubectl -n vault exec -it vault-0 -- vault operator init

2. Unseal Vault (run 3 times with different keys):
   kubectl -n vault exec -it vault-0 -- vault operator unseal

3. Configure Vault:
   - Enable Kubernetes authentication
   - Create policies
   - Configure secret engines

4. Deploy the application:
   kubectl apply -f k8s/app/

For more information, check the README.md
"
