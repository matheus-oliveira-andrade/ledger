#!/bin/bash

./install-nginx.sh

echo "Applying all manifests"
kubectl apply -f . --recursive