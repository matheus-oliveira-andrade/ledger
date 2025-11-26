#!/bin/bash

echo "Installing gateway api controller"
kubectl apply --server-side -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.4.0/standard-install.yaml

echo "Installing nginx gateway fabric"
helm install ngf oci://ghcr.io/nginx/charts/nginx-gateway-fabric --create-namespace -n nginx-gateway

echo "Waiting resources creation"
kubectl wait --timeout=5m -n nginx-gateway deployment/ngf-nginx-gateway-fabric --for=condition=Available

echo "Applying all manifests"
kubectl apply -f . --recursive