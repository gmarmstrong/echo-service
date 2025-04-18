# echo-service

Tiny Go HTTP server that simply echoes the request it receives. May be useful for debugging clients, load‑balancers, or service meshes.

![CI](https://github.com/gmarmstrong/echo-service/actions/workflows/build-image.yaml/badge.svg)

## Background reading

- <https://www.alexedwards.net/blog/an-introduction-to-handlers-and-servemuxes-in-go>

## TODOs

- Enforce a minimum test coverage
- More logging in ./cmd/echo-service/main.go
- Integration tests (run server and test the endpoints)
- End-to-end Kubernetes tests

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

## Endpoints

| Path       | Method | Purpose                   | Typical response        |
|------------|--------|---------------------------|-------------------------|
| `/{any}`   | GET    | Echo request back as JSON | 200 `application/json`  |
| `/healthz` | GET    | Liveness/readiness probe  | 204 No Content          |

The service listens on `PORT` (default **8080**).  
You can override it by setting that environment variable using kubectl or
by editing `k8s/deployment.yaml`:

```yaml
env:
  - name: PORT
    value: "9090"
ports:
  - containerPort: 9090
readinessProbe:
  httpGet:
    path: /healthz
    port: 9090
livenessProbe:
  httpGet:
    path: /healthz
    port: 9090
```

## Development & CI

### CI flow

The *Build container image* workflow (`.github/workflows/build-image.yml`) runs on every push to `main` (and on tag pushes). It:

1. Builds a multi‑arch Docker image with Buildx  
2. Tags it based on the Git ref (`v0.x.y`, `sha‑<digest>`)  
3. Pushes the image to GHCR

### Local testing

```sh
# Static analysis
go vet ./...

# Unit tests
go test ./... -v

# Check that the image builds and runs
docker build -t echo-service:dev .
docker run --rm -p 8080:8080 echo-service:dev
```

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

## License

This project uses the MIT license. See [LICENSE](LICENSE) for the full text.
