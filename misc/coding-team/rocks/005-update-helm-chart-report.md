# Completion Report: Update Helm chart for single rock image

## Summary
- Consolidated `images.cni` and `images.init` into a single `images.fancni` entry in `values.yaml`, pointing to `ghcr.io/ktsakalozos-canonical/fancni`.
- Updated `daemonset.yaml`: both `node-init` (init container) and `install-cni` (main container) now reference `images.fancni` and have explicit `command` overrides (`/init-node.sh` and `/install-cni.sh` respectively).
- `helm lint` and `helm template` pass with no errors.

## Files Changed
- `deploy/helm/fancni/values.yaml`
- `deploy/helm/fancni/templates/daemonset.yaml`

## Notable Tradeoffs / Risks
None. Changes are minimal and confined to the specified files. The `networkPolicy.image` (kube-router) was not touched.
