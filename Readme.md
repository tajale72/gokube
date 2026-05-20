# Gokube Weather API

A production-style weather forecasting service built in Go that integrates with the National Weather Service (NWS) API.

The project demonstrates backend engineering fundamentals including:

- REST API development in Go
- Middleware architecture
- API authentication
- Swagger/OpenAPI documentation
- Docker and Kubernetes deployment
- Minikube local infrastructure
- Unit testing with high coverage
- External API integration
- TLS-enabled HTTP server configuration
- Clean and maintainable Go project structure

---

# Architecture Overview

```text
Client → Middleware → Handlers → NWS API
                ↓
         Authentication
         Logging
         CORS
```

The service exposes weather forecast APIs while handling:

- Authentication
- Validation
- Retry-ready HTTP architecture
- Error handling
- Middleware processing
- Static frontend serving

---

# Features

- Weather forecast endpoint
- API key authentication middleware
- Health check endpoint
- Swagger/OpenAPI documentation
- Static frontend support
- Docker support
- Kubernetes manifests
- Minikube deployment
- Structured unit tests
- TLS configuration
- Middleware-based request processing

---

# Project Structure

```text
.
├── bin/
│   └── gokube
│
├── cmd/
│   └── gokube/
│       ├── main.go
│       └── main_test.go
│
├── config/
│   ├── config.go
│   └── config_test.go
│
├── router/
│   ├── router.go
│   ├── weather.go
│   ├── types.go
│   ├── router_test.go
│   └── weather_test.go
│
├── server/
│   ├── server.go
│   └── server_test.go
│
├── static/
│   ├── index.html
│   ├── index.css
│   ├── index.js
│   ├── swagger.html
│   └── swagger.yaml
│
├── k8s/
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   └── kustomization.yaml
│
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
├── server.crt
├── server.key
├── coverage.out
├── coverage.html
└── README.md
```

---

# API Endpoints

## Health Check

```http
GET /health
```

Response:

```text
go1.xx.x v1.0.0
```

---

## Weather Forecast

```http
GET /weather?lat=32.7767&long=-96.7970
```

### Headers

```http
X-API-KEY: 111
```

### Example Response

```json
{
  "short_forecast": "Sunny",
  "temperature": 90,
  "temperature_unit": "F",
  "characterization": "hot"
}
```

---

# Authentication

Protected endpoints require an API key header:

```http
X-API-KEY: 111
```

Authentication is enforced through middleware before requests reach the handler layer.

---

# Swagger Documentation

Swagger YAML:

```text
http://localhost:8080/swagger.yaml
```

Swagger UI:

```text
http://localhost:8080/docs
```

---

# Environment Variables

| Variable | Default |
|---|---|
| APP_ENV | local |
| MONGODB_URI | mongodb://localhost:27017 |
| AWS_REGION | us-east-1 |
| PORT | 8080 |
| SVC_NAME | my-service |
| WEATHER_API | https://api.weather.gov |
| X_API_KEY | default |

---

# Running Locally

## Install Dependencies

```bash
go mod tidy
```

---

## Run Application

```bash
make run
```

or directly:

```bash
go run ./cmd/gokube
```

Application runs on:

```text
http://localhost:8080
```

---

# Makefile Commands

## Show Available Commands

```bash
make help
```

---

## Build Application

```bash
make build
```

---

## Run Application

```bash
make run
```

---

## Build and Run

```bash
make all
```

---

## Clean Build Artifacts

```bash
make clean
```

---

# Docker

## Build Docker Image

```bash
make docker-build
```

---

## Run Docker Container

```bash
make docker-run
```

---

## Stop Docker Container

```bash
make docker-stop
```

---

# Kubernetes / Minikube

## Start Minikube

```bash
make minikube-start
```

---

## Deploy Application

```bash
make minikube-deploy
```

---

## Configure Local Host

Add the following to `/etc/hosts`:

```bash
<minikube-ip> gokube.local
```

Example:

```bash
192.168.49.2 gokube.local
```

---

## Access Application

```text
http://gokube.local
```

---

## Start Minikube Tunnel

```bash
make minikube-tunnel
```

---

## Delete Kubernetes Resources

```bash
make k8s-delete
```

---

## Stop Minikube

```bash
make minikube-stop
```

---

# Running Tests

## Execute Unit Tests

```bash
go test ./...
```

---

## Generate Coverage Report

```bash
go test ./... -coverprofile=coverage.out
```

---

## Generate HTML Coverage

```bash
go tool cover -html=coverage.out -o coverage.html
```

---

## Open Coverage Report

```bash
open coverage.html
```

---

# Frontend

Static frontend files are served from the `static/` directory.

Frontend includes:

- HTML
- CSS
- JavaScript
- Dynamic weather forecast fetching
- Swagger UI

---

# Middleware

## Logger Middleware

Logs request duration and request lifecycle information.

---

## Authentication Middleware

Protects `/weather` endpoint using `X-API-KEY`.

---

## CORS Middleware

Supports:

- GET
- POST
- PUT
- PATCH
- DELETE
- OPTIONS

---

# Security Features

- TLS 1.3 minimum version
- Restricted cipher suites
- API key authentication
- Request validation
- CORS middleware
- Protected weather endpoint

---

# Testing Coverage

The project includes unit tests for:

- Config package
- Middleware
- HTTP handlers
- Weather service logic
- Validation logic
- TLS server configuration
- Swagger handlers
- External API failure scenarios
- JSON encoding failures

Coverage includes:

- Success paths
- Failure paths
- Edge cases
- External API failures

---

# Technologies Used

- Golang
- net/http
- httptest
- Docker
- Kubernetes
- Minikube
- TLS
- Swagger/OpenAPI
- National Weather Service API

---

# Example Requests

## Health Check

```bash
curl http://localhost:8080/health
```

---

## Weather Forecast

```bash
curl --location 'http://localhost:8080/weather?lat=32.7767&long=-96.7970' \
--header 'X-API-KEY: 111'
```

---

# Future Improvements

- Structured logging
- Retry backoff strategy
- Rate limiting
- Metrics and observability
- Prometheus integration
- Distributed tracing
- CI/CD pipeline
- Helm chart support
- JWT authentication
- OpenTelemetry integration

---

# Author

Romit Tajale