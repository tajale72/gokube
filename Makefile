# Application Name
APP_NAME=gokube
DOCKER_IMAGE=$(APP_NAME)
DOCKER_USERNAME?=your-docker-username
K8S_NAMESPACE=$(APP_NAME)

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOSRC=$(GOBASE)/cmd/$(APP_NAME)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: all build clean run docker-build docker-run docker-stop help k8s-apply k8s-delete k8s-deploy minikube-start minikube-deploy minikube-tunnel minikube-stop

## Build and run the application
all: build run

## Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(GOBIN)
	go build -o $(GOBIN)/$(APP_NAME) ./cmd/$(APP_NAME)

## Run the application
run:
	@echo "Running $(APP_NAME)..."
	@if [ -f $(GOBIN)/$(APP_NAME) ]; then \
		$(GOBIN)/$(APP_NAME); \
	else \
		go run ./cmd/$(APP_NAME); \
	fi

## Clean build files
clean:
	@echo "Cleaning build files..."
	@rm -rf $(GOBIN)
	go clean

## Build docker image with Minikube's Docker daemon
docker-build:
	@echo "Building docker image using Minikube's Docker daemon..."
	@eval $$(minikube docker-env) && docker build -t $(DOCKER_IMAGE):latest .

## Run docker container
docker-run:
	@echo "Running docker container..."
	docker run --name $(APP_NAME) $(DOCKER_IMAGE):latest

## Stop and remove docker container
docker-stop:
	@echo "Stopping docker container..."
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

## Start Minikube cluster
minikube-start:
	@echo "Starting Minikube cluster..."
	minikube start
	minikube addons enable ingress
	@echo "Minikube IP: $$(minikube ip)"

## Deploy to Minikube
minikube-deploy: docker-build
	@echo "Deploying to Minikube..."
	kubectl create namespace $(K8S_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -k k8s/
	@echo "Waiting for deployment to be ready..."
	kubectl wait --namespace $(K8S_NAMESPACE) \
		--for=condition=ready pod \
		--selector=app=$(APP_NAME) \
		--timeout=90s
	@echo "Application URL: http://gokube.local"
	@echo "Add this to your /etc/hosts file:"
	@echo "$$(minikube ip) gokube.local"

## Start Minikube tunnel for LoadBalancer services
minikube-tunnel:
	@echo "Starting Minikube tunnel..."
	minikube tunnel

## Stop Minikube cluster
minikube-stop:
	@echo "Stopping Minikube cluster..."
	minikube stop

## Delete from Kubernetes
k8s-delete:
	@echo "Deleting from Kubernetes..."
	kubectl delete -k k8s/

## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  %-20s %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help 
