# Helm Chart Setup Guide

## Quick Start

### 1. Create Image Pull Secret for Artifactory

Before deploying, create a Kubernetes secret to pull images from Artifactory:

```bash
kubectl create secret docker-registry artifactory-registry-secret \
  --docker-server=YOUR_ARTIFACTORY_REGISTRY_URL \
  --docker-username=YOUR_USERNAME \
  --docker-password=YOUR_PASSWORD \
  --docker-email=YOUR_EMAIL@example.com \
  --namespace=gokube
```

Or if you already have the secret, reference it in your deployment.

### 2. Update values.yaml

Edit `helm/values.yaml` and set:
- `global.imageRegistry`: Your Artifactory registry URL
- `global.imagePullSecrets`: List of secret names (e.g., `[{"name": "artifactory-registry-secret"}]`)

### 3. Deploy with Helm

```bash
helm install gokube ./helm \
  --namespace gokube \
  --create-namespace \
  --set global.imageRegistry=your-artifactory.com/docker \
  --set api.image.tag=latest \
  --set ui.image.tag=latest
```

## GitLab CI/CD Setup

### Required Variables

In GitLab, go to **Settings > CI/CD > Variables** and add:

1. **ARTIFACTORY_REGISTRY**: Your Artifactory Docker registry URL
   - Example: `your-artifactory.com/docker`

2. **ARTIFACTORY_USER**: Your Artifactory username
   - Or use GitLab's built-in `CI_REGISTRY_USER` if available

3. **ARTIFACTORY_PASSWORD**: Your Artifactory password/token
   - Or use GitLab's built-in `CI_REGISTRY_PASSWORD` if available
   - Mark as **Protected** and **Masked**

4. **KUBECONFIG**: Base64 encoded kubeconfig file
   - Get it with: `cat ~/.kube/config | base64 -w 0`
   - Mark as **Protected** and **Masked**

5. **ARTIFACTORY_IMAGE_PULL_SECRET**: Name of the image pull secret in Kubernetes
   - Example: `artifactory-registry-secret`

### Alternative: GitLab Kubernetes Integration

Instead of using KUBECONFIG variable:

1. Go to **Operations > Kubernetes**
2. Click **Add Kubernetes cluster**
3. Follow the setup wizard
4. The CI/CD pipeline will automatically use cluster credentials

### Runner Tags

Ensure your GitLab runners have the appropriate tags:
- `docker`: For build and push jobs
- `kubernetes`: For deploy jobs

## Testing the Pipeline

1. Push code to `main` or `develop` branch
2. The pipeline will:
   - Build Docker images for API and UI
   - Push images to Artifactory
   - Deploy to Kubernetes using Helm

3. Monitor the pipeline in GitLab CI/CD > Pipelines

## Manual Deployment

If you need to deploy manually:

```bash
# Build and push images manually
docker build -t your-artifactory.com/docker/ccda-api:latest -f Dockerfile .
docker build -t your-artifactory.com/docker/vue-ui:latest -f Dockerfile.ui .
docker push your-artifactory.com/docker/ccda-api:latest
docker push your-artifactory.com/docker/vue-ui:latest

# Deploy with Helm
helm upgrade --install gokube ./helm \
  --namespace gokube \
  --create-namespace \
  --set global.imageRegistry=your-artifactory.com/docker \
  --set api.image.tag=latest \
  --set ui.image.tag=latest
```

## Troubleshooting

### Images not pulling from Artifactory

1. Verify the image pull secret exists:
   ```bash
   kubectl get secret artifactory-registry-secret -n gokube
   ```

2. Check if the secret is referenced in the deployment:
   ```bash
   kubectl describe deployment ccda-api -n gokube | grep -A 5 "Image Pull Secrets"
   ```

3. Verify image exists in Artifactory and is accessible

### Helm deployment fails

1. Check Helm release status:
   ```bash
   helm status gokube -n gokube
   ```

2. View Helm release history:
   ```bash
   helm history gokube -n gokube
   ```

3. Debug Helm template rendering:
   ```bash
   helm template gokube ./helm --debug
   ```

### Pods not starting

1. Check pod status:
   ```bash
   kubectl get pods -n gokube
   ```

2. View pod events:
   ```bash
   kubectl describe pod <pod-name> -n gokube
   ```

3. Check pod logs:
   ```bash
   kubectl logs <pod-name> -n gokube
   ```

## Updating the Deployment

### Update Image Tags

```bash
helm upgrade gokube ./helm \
  --namespace gokube \
  --set api.image.tag=v1.1.0 \
  --set ui.image.tag=v1.1.0
```

### Scale Replicas

```bash
helm upgrade gokube ./helm \
  --namespace gokube \
  --set api.replicaCount=3 \
  --set ui.replicaCount=3
```

### Rollback

```bash
helm rollback gokube -n gokube
```

Or use the manual rollback job in GitLab CI/CD pipeline.

