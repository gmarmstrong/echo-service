# echo-service

Tiny Go HTTP server that simply echoes the request it receives. May be useful for debugging clients, load‑balancers, or service meshes.

![CI](https://github.com/gmarmstrong/echo-service/actions/workflows/build-image.yml/badge.svg)

## Usage

### Basic usage (with Go, no Docker)

```sh
go run ./cmd/echo-service
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

- TODO

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
