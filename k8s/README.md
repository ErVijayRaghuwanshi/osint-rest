# Kubernetes Deployment Guide – OSINT Scraper

This directory contains all **Kubernetes manifests** required to deploy the **OSINT Scraper** service in a Kubernetes cluster.

The setup is optimized for **local development using Minikube**, while remaining **production-friendly** with minimal changes.

---

## Directory Structure

```
k8s/
├── namespace.yaml      # Namespace isolation
├── configmap.yaml      # Application configuration
├── deployment.yaml     # Application deployment
├── service.yaml        # Internal service exposure
├── ingress.yaml        # Optional external access
└── README.md           # This file
```

---

## Prerequisites

* Kubernetes cluster

  * **Minikube** (recommended for local development)
* kubectl
* Docker
* Makefile (recommended)

Verify installations:

```bash
kubectl version --client
minikube version
docker version
```

---

## Recommended Local Setup (Minikube)

### 1️⃣ Start Minikube

```bash
minikube start
```

Ensure kubectl context is correct:

```bash
kubectl config current-context
```

---

## Building Docker Image for Kubernetes

⚠️ **Important**: Minikube uses its **own Docker runtime**. Images built on the host Docker daemon are **not visible** to Kubernetes.

### Correct way (using Makefile)

```bash
echo "0.2.0" > version
make docker-build-k8s
```

This builds the image **inside Minikube**:

```text
osint-scraper:0.2.0
```

---

## Deploying to Kubernetes

### Apply all manifests

```bash
make k8s-apply
```

This will create:

* Namespace (`osint`)
* ConfigMap
* Deployment (2 replicas)
* Service
* Ingress (if enabled)

---

## Verifying Deployment

```bash
kubectl get all -n osint
```

Check logs:

```bash
kubectl logs -n osint deploy/osint-scraper
```

Restart deployment after image update:

```bash
make k8s-restart
```

---

## Accessing the Service

### Option 1: Minikube Service (Recommended)

```bash
minikube service osint-scraper -n osint
```

This will open the service URL in your browser.

---

### Option 2: Ingress (Optional)

If you enabled `ingress.yaml` and have an ingress controller installed:

```bash
minikube addons enable ingress
```

Add to `/etc/hosts`:

```text
127.0.0.1 osint.local
```

Access:

```text
http://osint.local
```

---

## Configuration

Application configuration is managed via **ConfigMap**:

```yaml
PORT: "8080"
LOG_LEVEL: "info"
```

To update configuration:

```bash
kubectl apply -f k8s/configmap.yaml
kubectl rollout restart deployment/osint-scraper -n osint
```

---

## Image Versioning Strategy

Best practice:

1. Update version

   ```bash
   echo "0.2.1" > version
   ```
2. Build image inside Minikube

   ```bash
   make docker-build-k8s
   ```
3. Restart deployment

   ```bash
   make k8s-restart
   ```

---

## Production Notes

For production clusters:

* Push image to a container registry (Docker Hub / ECR / GCR)
* Update `deployment.yaml`:

  ```yaml
  image: <registry>/osint-scraper:0.2.0
  imagePullPolicy: Always
  ```
* Use **Secrets** instead of ConfigMaps for sensitive data
* Enable:

  * HPA (Horizontal Pod Autoscaler)
  * NetworkPolicies
  * PodSecurityStandards

---

## Cleanup

Remove all OSINT resources:

```bash
make k8s-delete
```

---

## Troubleshooting

### ImagePullBackOff

* Ensure image was built using:

  ```bash
  make docker-build-k8s
  ```
* Ensure `imagePullPolicy: IfNotPresent`

### Pods not ready

```bash
kubectl describe pod -n osint <pod-name>
```

---

## Maintainer

**Er Vijay Raghuwanshi**
OSINT • Backend • Distributed Systems

---

> This Kubernetes setup is intentionally simple, secure, and extensible — suitable for OSINT pipelines, scraping workloads, and future scaling.
