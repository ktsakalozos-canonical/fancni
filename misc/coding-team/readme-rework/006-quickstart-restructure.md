# Task: Restructure Quick Start in README.md

## Context

The README.md Quick Start currently shows build steps first. We want it to show the end-user deployment flow instead.

## Objective

Rewrite the "Quick Start" section of `README.md` to follow this sequence:

1. **Deploy Canonical K8s without the network** — show the bootstrap config and command:
   ```bash
   cat <<EOF > bootstrap-config.yaml
   cluster-config:
     network:
       enabled: false
   EOF
   sudo k8s bootstrap --file bootstrap-config.yaml
   ```

2. **Install fancni with Helm** — the existing `helm install` command:
   ```bash
   helm install fancni deploy/helm/fancni/
   ```

3. **Verify** — the existing verification step (kubectl get pods).

4. **Connectivity check** — move the content from the current "Connectivity Test" section here (the `kubectl apply`, get pods, ping commands).

## Scope

- Only edit `README.md`.
- Remove the old "Connectivity Test" section (its content moves into Quick Start step 4).
- Remove the old "Build container images" and "Load images onto cluster nodes" steps from Quick Start (those belong in the Development doc, and are already there).

## Non-goals

- Don't touch any other files.
- Don't change anything outside Quick Start and the Connectivity Test section.
- Don't rewrite prose in other sections.

## Constraints

- Keep the same heading level structure (###) for sub-steps within Quick Start.
- Keep the step numbering as 1, 2, 3, 4.
