# 1. Start minikube
minikube start
minikube addons enable ingress

# 2. Build image in minikube context
eval $(minikube docker-env)
docker build -t gokube:latest .

# 3. Deploy to Kubernetes
kubectl apply -f k8s/

# 4. Check deployment
kubectl get pods -l app=gokube
kubectl get svc gokube