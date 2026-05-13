# Task: Rewrite README.md

## Context
We've extracted development and configuration docs. Now rewrite the README for a contributor audience.

## Objective
Rewrite `/home/jackal/workspace/fancni/README.md` with the following structure:

### Structure (in order)

1. **Title + one-liner** — keep the existing description
2. **Alpha notice** — short callout (e.g., `> ⚠️ **Alpha** — not yet production-ready.`)
3. **How It Works** — move the existing "How It Works" paragraph here (explains Fan networking mechanism). Keep it concise.
4. **Prerequisites**:
   - Ubuntu nodes (Fan networking relies on Ubuntu kernel support)
   - Canonical Kubernetes cluster (recommend it; link to https://ubuntu.com/kubernetes)
   - Helm 3
   - Docker (for building images)
   - **No mention of k3s anywhere**
5. **Quick Start**:
   - Build: `make docker-build`
   - Load images onto nodes — show only the containerd/ctr approach for Canonical Kubernetes:
     ```
     docker save fancni:latest | sudo ctr -n k8s.io images import -
     docker save fancni-init:latest | sudo ctr -n k8s.io images import -
     ```
   - Helm install: `helm install fancni deploy/helm/fancni/`
   - Verify: `kubectl -n kube-system get pods -l app=fancni`
6. **Architecture** — single line linking to `ARCHITECTURE.md`
7. **Configuration** — single line linking to `docs/configuration.md`
8. **Connectivity Test** — keep existing section (lines 79-93 of current README), remove k3s references if any
9. **Troubleshooting** — keep existing section as-is
10. **Development** — single line linking to `docs/development.md`

## Constraints
- Zero references to k3s
- Contributor-focused tone (assume reader knows Kubernetes basics)
- Keep it scannable — no walls of text

## Non-goals
- Don't touch any other files
- Don't add comparison tables or marketing copy
