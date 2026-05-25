# Inventory Manager

A containerized REST API built with GoLang for managing inventory items. Deployed via Docker, orchestrated with Kubernetes (minikube), and automated through a GitHub Actions CI/CD pipeline.

## Tech Stack

| Layer | Technology |
| --- | --- |
| Language | Go |
| Containerization | Docker |
| Orchestration | Kubernetes (minikube) |
| CI/CD | GitHub Actions |
| Local Dev | docker-compose |

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Docker](https://www.docker.com/products/docker-desktop/)
- [minikube](https://minikube.sigs.k8s.io/docs/start/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [make](https://www.gnu.org/software/make/)

## Project Structure

```text
.
├── .github/workflows/
│   ├── ci.yml              # Lint, test, build on every push/PR
│   └── cd.yml              # Build & push image, deploy to k8s
├── cmd/server/
│   └── main.go             # Application entrypoint
├── internal/
│   ├── handlers/           # HTTP route handlers
│   ├── middleware/         # HTTP middleware (logging, etc.)
│   ├── models/             # Data structs
│   └── repository/         # Data access layer
├── k8s/
│   ├── namespace.yaml
│   ├── configmap.yaml
│   ├── deployment.yaml
│   └── service.yaml
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── Makefile
└── .env.example            # Environment variable template
```

## Getting Started

### 1. Environment

Copy the example env file and fill in your values:

```bash
cp .env.example .env
```

### 2. Run Locally (without Docker)

```bash
make run
```

### 3. Run with Docker Compose

```bash
make docker-up
```

### 4. Deploy to Kubernetes (minikube)

```bash
minikube start
make k8s-deploy
```

## API Endpoints

| Method | Path | Description |
| --- | --- | --- |
| GET | `/health` | Health check |
| GET | `/items` | List all items |
| GET | `/items/{id}` | Get item by ID |
| POST | `/items` | Create a new item |
| PUT | `/items/{id}` | Update an item |
| DELETE | `/items/{id}` | Delete an item |

## CI/CD Pipeline

**CI (`ci.yml`)** — triggers on every push and pull request:

1. Lint with `golangci-lint`
2. Run unit tests
3. Build the binary

**CD (`cd.yml`)** — triggers on merge to `main`:

1. Build and push Docker image to registry
2. Update the Kubernetes deployment image tag
3. Commit updated manifest back to repo (GitOps — run `make k8s-deploy` locally to apply)

## Makefile Targets

```bash
make run          # Run the server locally
make test         # Run tests
make build        # Build the binary
make docker-up    # Start services with docker-compose
make docker-down  # Stop docker-compose services
make k8s-deploy   # Apply all k8s manifests
make k8s-delete   # Tear down k8s resources
```

## Environment Variables

| Variable | Description | Default |
| --- | --- | --- |
| `PORT` | Port the server listens on | `8080` |
| `ENV` | Runtime environment | `development` |
| `DB_URL` | Database connection string | — |
