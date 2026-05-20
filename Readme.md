# GoKube - Resident Portal Application

A production-grade Go web application with a resident portal frontend, designed to run in Kubernetes environments.

## 🏗️ Architecture

- **Backend**: Go HTTP server with TLS support
- **Frontend**: HTML/CSS/JavaScript resident portal
- **Database**: PostgreSQL with GORM ORM
- **Deployment**: Docker containers with Kubernetes manifests
- **Security**: TLS 1.3 with modern cipher suites

## 🚀 Features

- **Resident Portal**: Modern web interface for apartment residents
- **Health Checks**: Built-in health monitoring endpoints
- **TLS Security**: Production-ready HTTPS configuration
- **Kubernetes Ready**: Complete K8s deployment manifests
- **Static File Serving**: Efficient static asset delivery
- **Logging**: Structured logging with request timing

## 📋 Prerequisites

- Go 1.24.3 or later
- Docker
- Kubernetes cluster (minikube, kind, or cloud)
- PostgreSQL (for database features)

## 🛠️ Local Development

### 1. Clone and Setup

```bash
git clone <repository-url>
cd gokube
go mod download
```

### 2. Generate TLS Certificates

```bash
# Generate self-signed certificates for local development
openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes \
  -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"
```

### 3. Run Locally

```bash
# Build and run
go run cmd/gokube/main.go

# Or build binary
go build -o bin/gokube cmd/gokube/main.go
./bin/gokube
```

The application will be available at: `https://localhost:8080`

## 🐳 Docker Deployment

### 1. Build Docker Image

```bash
docker build -t gokube:latest .
```

### 2. Run Container

```bash
docker run -p 8080:8080 gokube:latest
```

## ☸️ Kubernetes Deployment (Minikube)

### 1. Start Minikube

```bash
minikube start
```

### 2. Build Image in Minikube

```bash
# Build image directly in minikube
eval $(minikube docker-env)
docker build -t gokube:latest .
```

### 3. Deploy to Kubernetes

```bash
# Apply all Kubernetes manifests
kubectl apply -f k8s/

# Or apply individually
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml
```

### 4. Access the Application

```bash
# Get the service URL
minikube service gokube --url

# Or if using ingress
echo "$(minikube ip) gokube.local" | sudo tee -a /etc/hosts
# Then visit: http://gokube.local
```

## 📁 Project Structure

```
gokube/
├── cmd/gokube/          # Main application entry point
├── config/              # Configuration management
├── router/              # HTTP routing and handlers
├── server/              # HTTP server configuration
├── static/              # Frontend assets (HTML, CSS, JS)
├── k8s/                 # Kubernetes manifests
├── Dockerfile           # Container build instructions
├── docker-compose.yml   # Local development setup
└── Makefile            # Build automation
```

## 🔧 Configuration

The application uses environment variables for configuration:

- `ADDR`: Server address (default: localhost)
- `PORT`: Server port (default: 8080)
- `MY_CERT`: TLS certificate file (default: server.crt)
- `MY_KEY`: TLS private key file (default: server.key)

## 🏥 Health Checks

- **Health Endpoint**: `GET /health`
- **Root Endpoint**: `GET /` (serves resident portal)

## 🎨 Resident Portal Features

- **Community Announcements**: Real-time updates
- **Quick Actions**: Maintenance requests, rent payment, profile access
- **Responsive Design**: Mobile-friendly interface
- **Modern UI**: Clean, professional appearance

## 🔒 Security Features

- **TLS 1.3**: Latest encryption standards
- **Modern Cipher Suites**: AES-256-GCM, ChaCha20-Poly1305
- **Secure Headers**: Production-ready security configuration
- **Input Validation**: Request sanitization

## 📊 Monitoring

The application includes:
- Request timing logs
- Health check endpoints
- Structured logging
- Kubernetes readiness/liveness probes

## 🚀 Production Deployment

For production deployment:

1. **Use proper TLS certificates** from a trusted CA
2. **Configure environment variables** for your infrastructure
3. **Set up database connections** for persistent data
4. **Configure ingress** with proper domain names
5. **Set resource limits** based on your requirements

## 🛠️ Development Commands

```bash
# Run tests
go test ./...

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o bin/gokube-linux cmd/gokube/main.go

# Format code
go fmt ./...

# Lint code
golangci-lint run
```

## 📝 API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Resident portal homepage |
| GET | `/health` | Health check endpoint |
| GET | `/static/*` | Static file serving |

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Troubleshooting

### Common Issues

1. **TLS Certificate Errors**: Ensure certificates are in the project root
2. **Port Conflicts**: Change the PORT environment variable
3. **Kubernetes Image Pull Errors**: Build image in minikube context
4. **Static Files Not Loading**: Check file paths and permissions

### Getting Help

- Check the logs: `kubectl logs -f deployment/gokube`
- Verify service: `kubectl get svc gokube`
- Check pods: `kubectl get pods -l app=gokube`