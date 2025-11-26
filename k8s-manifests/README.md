## Kubernetes manisfests for ledger application
This folder contains all required Kubernetes manifests (YAML files) to deploy and run the application on a Kubernetes cluster.

### Prerequisites
Ensure the following tools are installed:
- [Minikube](https://minikube.sigs.k8s.io/docs/start/) – for local Kubernetes cluster
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) – Kubernetes command-line tool
- [helm](https://helm.sh/docs/intro/install) – package manager for k8s

### Creating resources

- Creating cluster
```bash
minikube start
```

- Apply manifests on cluster
```bash
./up.sh
```

### Running the Application Locally

1 - Start the minikube tunnel in a new terminal
```bash    
minikube tunnel 
```

2 - Retrieve IP
```bash
kubectl get service nginx-gateway-nginx -o jsonpath='{.spec.clusterIP}'
# 10.108.78.123
```

*Obs.: all HTTP requests can be found [here](../docs/assets/ledger%20http%20requests.postman_collection.json)*

### Exposing database and GRPC server locally

In a new terminal
```sh
kubectl port-forward svc/account-db-svc 5432:5432
# kubectl port-forward svc/ledger-db-svc 5433:5433
# kubectl port-forward svc/account-api-service 9001:9001 # account GRPC, accessible via localhost
```

### Deleting resouces
To remove all resources created by the manifests

- Delete all resources applied on cluster
```bash
./down.sh
```

- Delete Kubernetes cluster
```bash
minikube delete
```
