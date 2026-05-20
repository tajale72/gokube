# Dragon Weather API

A lightweight Go weather service that fetches live forecast data from the National Weather Service (NWS) API.

This project demonstrates:

- REST API development in Go
- Middleware architecture
- API authentication
- CORS support
- Unit testing
- Error handling
- External API integration
- Clean project structure

---

# Features

- Weather forecast endpoint
- API key authentication
- CORS middleware
- Health check endpoint
- Static frontend support
- 100% unit test coverage
---

# Project Structure

```text
.
├── bin
│   └── gokube
├── cmd
│   └── gokube
│       ├── main_test.go
│       └── main.go
├── config
│   ├── config_test.go
│   └── config.go
├── coverage.html
├── coverage.out
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── k8s
│   ├── deployment.yaml
│   ├── ingress.yaml
│   ├── kustomization.yaml
│   └── service.yaml
├── Makefile
├── Readme.md
├── router
│   ├── router_test.go
│   ├── router.go
│   ├── types.go
│   ├── weather_test.go
│   └── weather.go
├── server
│   ├── server_test.go
│   └── server.go
├── server.crt
├── server.key
├── static
│   ├── index.css
│   ├── index.html
│   └── index.js
└── swagger.yaml