# Task 006: Helm Chart

## Context

All code, containers, and scripts are built. This task creates the Helm chart that deploys fancni as a DaemonSet on every node.

## Objective

Create a Helm chart under `deploy/helm/fancni/` that deploys:
- A DaemonSet with two init containers (node-init + cni-install) and a pause container
- A ConfigMap for the CNI network config
- A ServiceAccount (minimal, no RBAC needed since we don't talk to the K8s API)

## Scope

### deploy/helm/fancni/Chart.yaml

```yaml
apiVersion: v2
name: fancni
description: Fan CNI - Kubernetes CNI plugin using Ubuntu Fan Networking
type: application
version: 0.1.0
appVersion: "0.1.0"
```

### deploy/helm/fancni/values.yaml

```yaml
overlay:
  network: "240.0.0.0/8"

underlay:
  prefix: 16

images:
  cni:
    repository: fancni
    tag: latest
    pullPolicy: IfNotPresent
  init:
    repository: fancni-init
    tag: latest
    pullPolicy: IfNotPresent

cni:
  version: "1.0.0"
  name: "fancni"
  confFileName: "10-fancni.conf"

ipam:
  dataDir: "/var/lib/cni/fancni"

tolerations:
  - operator: Exists
    effect: NoSchedule
  - operator: Exists
    effect: NoExecute
```

### deploy/helm/fancni/templates/configmap.yaml

Creates the CNI config JSON that gets written to `/etc/cni/net.d/`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name }}-cni-config
  namespace: {{ .Release.Namespace }}
data:
  cni-conf: |
    {
      "cniVersion": "{{ .Values.cni.version }}",
      "name": "{{ .Values.cni.name }}",
      "type": "fancni",
      "overlayNetwork": "{{ .Values.overlay.network }}",
      "underlayPrefix": {{ .Values.underlay.prefix }}
    }
```

### deploy/helm/fancni/templates/daemonset.yaml

Key aspects:
- `hostNetwork: true` — the pod runs in the host network namespace
- `hostPID: true` — needed for `nsenter --target 1`
- Tolerations to run on all nodes including masters
- Node selector: only Ubuntu nodes (use a label or omit for now)

**Init container 1: node-init**
- Image: `{{ .Values.images.init.repository }}:{{ .Values.images.init.tag }}`
- Privileged: true
- Env: `OVERLAY_NETWORK`, `UNDERLAY_PREFIX` from values
- Volume mounts:
  - `/host` → host root `/` (for chroot)

**Init container 2: install-cni**  
- Image: `{{ .Values.images.cni.repository }}:{{ .Values.images.cni.tag }}`
- Env: `CNI_CONF` from ConfigMap
- Volume mounts:
  - `/host/opt/cni/bin` → hostPath `/opt/cni/bin`
  - `/host/etc/cni/net.d` → hostPath `/etc/cni/net.d`

Wait — `install-cni` should NOT be an init container since it runs `sleep infinity` to keep the pod alive. It should be the **main container**.

Corrected design:
- **Init container**: `node-init` (runs once, exits)
- **Main container**: `install-cni` (copies binary + config, then sleeps forever)

**Main container: install-cni**
- Image: `{{ .Values.images.cni.repository }}:{{ .Values.images.cni.tag }}`
- Env: `CNI_CONF` from ConfigMap
- Volume mounts:
  - `/host/opt/cni/bin` → hostPath `/opt/cni/bin`
  - `/host/etc/cni/net.d` → hostPath `/etc/cni/net.d`

### deploy/helm/fancni/templates/serviceaccount.yaml

Simple ServiceAccount (no RBAC bindings needed):

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Release.Name }}-fancni
  namespace: {{ .Release.Namespace }}
```

### Makefile update

Add:
```makefile
helm-template:
	helm template fancni deploy/helm/fancni/

helm-lint:
	helm lint deploy/helm/fancni/
```

## Non-goals

- No RBAC (ClusterRole/ClusterRoleBinding) — we don't access the K8s API
- No PodSecurityPolicy / PodSecurity admission
- No Ingress / Service
- No HPA

## Constraints

- DaemonSet pod must tolerate all taints (to run on control plane nodes too)
- The `node-init` init container must be privileged with `hostPID: true`
- The `install-cni` main container does NOT need to be privileged
- `helm lint` must pass
