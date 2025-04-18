# echo-service

Tiny Go HTTP server that simply echoes the request it receives. May be useful for debugging clients, load‑balancers, or service meshes.

![CI](https://github.com/gmarmstrong/echo-service/actions/workflows/build-image.yml/badge.svg)

## Requirements

Requirements vary depending on how you want to use this.

* jq for pretty-printing JSON output (optional)
* Go 1.22+ (only for `go run` / local builds)  
* Docker 20+ (only for the container demo)  
* kubectl and a cluster (for the Kubernetes section)
* Optionally, minikube (for the Kubernetes section)

## Usage

### Quick start (run from source)

```sh
go run ./cmd/echo-service
# In another terminal:
curl -s http://localhost:8080/hello | jq .
```

### Docker usage (with Docker, no K8s)

From the repository root, run:

```sh
docker build -t echo-service:0.1.0 .
docker run --rm -p 8080:8080 echo-service:0.1.0
```

Then in another terminal instance, run:

```sh
curl -s http://localhost:8080/hello | jq .
```

and you should get

```json
{
  "method": "GET",
  "path": "/hello",
  "headers": {
    "Accept": [
      "*/*"
    ],
    "User-Agent": [
      "curl/8.12.1"
    ]
  },
  "host": "localhost:8080",
  "remote_addr": "[::1]:65325"
}
```

### Kubernetes usage

The manifests in `k8s/` let you spin up the service on any cluster.  

These instructions use minikube.
On Kind or a real cluster, just skip the minikube‑specific commands.

```sh
# Start a local cluster
minikube start

# Build or load the image inside minikube's container runtime
eval "$(minikube -p minikube docker-env)"
docker build -t echo-service:0.1.0 .

# Deploy the manifests
kubectl apply -f k8s/
kubectl rollout status deployment/echo-service

# Access the service
minikube service echo-service --url       # prints http://<node-ip>:<port>
curl -s $(minikube service echo-service --url)/hello | jq .
```

If port 8080 is reserved on your machine, you can change the port in `k8s/deployment.yaml`
or use `kubectl set env deployment/echo-service PORT=9090`.

## Project structure

Kubernetes manifests are stored in the `k8s/` directory for convenience.
This works for our single-binary, demo use case.

```
# k8s/
# ├── deployment.yaml
# └── service.yaml
```

In production, you should use Kustomize, Helm, or similar, and move the
manifests to a dedicated GitOps repo. That will give you clean separation
of app code vs. cluster state.

Alternatively, if you prefer monorepos, you could do something like:

```
├── apps
│   ├── base
│   ├── production 
│   └── staging
├── infrastructure
│   ├── base
│   ├── production 
│   └── staging
└── clusters
    ├── production
    └── staging
```
