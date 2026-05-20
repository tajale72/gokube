# Gokube Helm Chart

This Helm chart deploys the CCDA-API service (Golang) and Vue UI application to Kubernetes.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- Access to Artifactory Docker registry
- kubectl configured to access your Kubernetes cluster

## Installation

### Basic Installation

```bash
helm install gokube ./helm \
  --namespace gokube \
  --create-namespace \
  --set global.imageRegistry=your-artifactory.com/docker \
  --set api.image.tag=v1.0.0 \
  --set ui.image.tag=v1.0.0
```

### Using Environment-Specific Values

For development:
```bash
helm install gokube ./helm \
  -f helm/values-dev.yaml \
  --namespace gokube \
  --create-namespace \
  --set global.imageRegistry=your-artifactory.com/docker
```

For production:
```bash
helm install gokube ./helm \
  -f helm/values-prod.yaml \
  --namespace gokube \
  --create-namespace \
  --set global.imageRegistry=your-artifactory.com/docker
```

## Configuration

### Image Pull Secrets

If your Artifactory requires authentication, create a Kubernetes secret:

```bash
kubectl create secret docker-registry artifactory-registry-secret \
  --docker-server=your-artifactory.com \
  --docker-username=your-username \
  --docker-password=your-password \
  --docker-email=your-email@example.com \
  --namespace gokube
```

Then reference it in your Helm values:
```yaml
global:
  imagePullSecrets:
    - name: artifactory-registry-secret
```

### Updating Values

Edit `values.yaml` or use `--set` flags:

```bash
helm upgrade gokube ./helm \
  --namespace gokube \
  --set api.replicaCount=3 \
  --set ui.replicaCount=2
```

## Upgrading

```bash
helm upgrade gokube ./helm \
  --namespace gokube \
  --set api.image.tag=v1.1.0 \
  --set ui.image.tag=v1.1.0
```

## Uninstalling

```bash
helm uninstall gokube --namespace gokube
```

## Values Reference

### Global Values

| Parameter | Description | Default |
|-----------|-------------|---------|
| `global.namespace` | Kubernetes namespace | `gokube` |
| `global.imageRegistry` | Artifactory registry URL | `""` |
| `global.imagePullSecrets` | Image pull secrets | `[]` |

### API Service Values

| Parameter | Description | Default |
|-----------|-------------|---------|
| `api.enabled` | Enable API deployment | `true` |
| `api.name` | API service name | `ccda-api` |
| `api.image.repository` | Image repository | `ccda-api` |
| `api.image.tag` | Image tag | `latest` |
| `api.replicaCount` | Number of replicas | `2` |
| `api.service.port` | Service port | `80` |
| `api.service.targetPort` | Container port | `8080` |

### UI Service Values

| Parameter | Description | Default |
|-----------|-------------|---------|
| `ui.enabled` | Enable UI deployment | `true` |
| `ui.name` | UI service name | `vue-ui` |
| `ui.image.repository` | Image repository | `vue-ui` |
| `ui.image.tag` | Image tag | `latest` |
| `ui.replicaCount` | Number of replicas | `2` |
| `ui.service.port` | Service port | `80` |
| `ui.service.targetPort` | Container port | `80` |

### Ingress Values

| Parameter | Description | Default |
|-----------|-------------|---------|
| `ingress.enabled` | Enable ingress | `true` |
| `ingress.className` | Ingress class name | `nginx` |
| `ingress.hosts` | Ingress hosts configuration | See values.yaml |

## GitLab CI/CD Integration

The `.gitlab-ci.yml` file includes:
- Building Docker images for API and UI
- Pushing images to Artifactory
- Deploying to Kubernetes using Helm

### Required GitLab CI/CD Variables

Set these in GitLab CI/CD Settings > Variables:

- `ARTIFACTORY_REGISTRY`: Your Artifactory Docker registry URL
- `ARTIFACTORY_USER`: Artifactory username (or use `CI_REGISTRY_USER`)
- `ARTIFACTORY_PASSWORD`: Artifactory password (or use `CI_REGISTRY_PASSWORD`)
- `KUBECONFIG`: Base64 encoded kubeconfig file (or use GitLab Kubernetes integration)
- `ARTIFACTORY_IMAGE_PULL_SECRET`: Name of the image pull secret in Kubernetes

### GitLab Kubernetes Integration

Alternatively, you can use GitLab's built-in Kubernetes integration:
1. Go to Operations > Kubernetes
2. Add your Kubernetes cluster
3. The CI/CD pipeline will automatically use the cluster credentials

## Troubleshooting

### Check Pod Status
```bash
kubectl get pods -n gokube
```

### View Pod Logs
```bash
kubectl logs -n gokube -l app.kubernetes.io/name=ccda-api
kubectl logs -n gokube -l app.kubernetes.io/name=vue-ui
```

### Check Helm Release
```bash
helm list -n gokube
helm status gokube -n gokube
```

### Debug Helm Template
```bash
helm template gokube ./helm --debug
```

