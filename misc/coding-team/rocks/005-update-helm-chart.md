# Task: Update Helm chart for single rock image

## Context
We now have a single rock image containing both the CNI binary/installer and the node-init scripts. The helm chart currently references two separate images (`images.cni` and `images.init`). We need to consolidate to one image reference and add command overrides.

## Objective
Update the Helm chart so both containers (init and main) use a single image with command overrides to select behavior.

## Scope
- `deploy/helm/fancni/values.yaml`
- `deploy/helm/fancni/templates/daemonset.yaml`

## Changes

### values.yaml
Replace:
```yaml
images:
  cni:
    repository: ghcr.io/ktsakalozos-canonical/fancni
    tag: latest
    pullPolicy: IfNotPresent
  init:
    repository: ghcr.io/ktsakalozos-canonical/fancni-init
    tag: latest
    pullPolicy: IfNotPresent
```
With:
```yaml
images:
  fancni:
    repository: ghcr.io/ktsakalozos-canonical/fancni
    tag: latest
    pullPolicy: IfNotPresent
```

### daemonset.yaml
- Init container (`node-init`): change image to `{{ .Values.images.fancni.repository }}:{{ .Values.images.fancni.tag }}`, add `command: ["/init-node.sh"]`, update `imagePullPolicy`
- Main container (`install-cni`): change image to `{{ .Values.images.fancni.repository }}:{{ .Values.images.fancni.tag }}`, add `command: ["/install-cni.sh"]`, update `imagePullPolicy`

## Non-goals
- Do not touch other templates (network-policy-daemonset, configmap, etc.)
- Do not change the networkPolicy.image reference (kube-router is a separate image)
