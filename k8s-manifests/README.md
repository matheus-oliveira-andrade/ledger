## Kubernetes manisfests for ledger application
This folder contains all required Kubernetes manifests (YAML files) to deploy and run the application on a Kubernetes cluster.

### Prerequisites
Ensure the following tools are installed:
- [Minikube](https://minikube.sigs.k8s.io/docs/start/) – for local Kubernetes cluster
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) – Kubernetes command-line tool


### Creating resources

- Creating cluster
```bash
minikube start
```

- Apply manifests
```bash
kubectl apply -Rf . 
```

### Running the Application Locally

By default, `minikube` does not expose cluster services to `localhost`. To access applications locally:

1 - Enable the ingress Add-on:
```bash
minikube addons enable ingress
```

2 - Start the inikube tunnel in a new terminal
```bash    
minikube tunnel 
```

Now you can access services exposed via the Ingress controller at `localhost` 

### Deleting resouces
To remove all resources created by the manifests

```bash
# delete all resources created with manifests file
kubectl delete -Rf . 
```


