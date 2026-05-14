# Task: Update e2e test to use rock instead of Docker

## Context
The e2e test (`tests/e2e/test-e2e.sh`) currently builds Docker images with `make docker-build`, transfers them with `docker save | ctr import`, and re-tags them. We now use a single rock (`.rock` file is an OCI archive) that can be imported directly with `ctr images import`.

## Objective
Replace the Docker-based image build and transfer flow with a rock-based flow.

## Scope
- `tests/e2e/test-e2e.sh`

## Changes needed

### Phase 6 — Replace docker build with rock build
Replace:
```bash
make -C "${REPO_ROOT}" docker-build
```
With:
```bash
make -C "${REPO_ROOT}" rock-build
```

The rock file will be at `${REPO_ROOT}/fancni_0.1_amd64.rock`.

### Phase 7 — Replace docker save/import with rock import
The `.rock` file is an OCI archive. Replace the docker save + ctr import loop with:

1. Push the `.rock` file into both VMs via `lxc file push`
2. Import with `ctr --address "${CTR_SOCK}" -n "${CTR_NS}" images import --all-platforms /tmp/fancni_0.1_amd64.rock`
3. Re-tag the imported image to match the Helm chart reference (`ghcr.io/ktsakalozos-canonical/fancni:latest`)

The rock's image name after import will be based on its OCI annotations. The image reference in the OCI archive is typically `fancni:0.1`. After import, re-tag it:
```bash
ctr --address "${CTR_SOCK}" -n "${CTR_NS}" images tag "fancni:0.1" "ghcr.io/ktsakalozos-canonical/fancni:latest"
```

4. Still pull and import `nginx:latest` and `cloudnativelabs/kube-router:v2.9.0` via docker (those are external images unrelated to the rock).
5. Remove the old fancni-init re-tag logic (there's only one image now).

### Summary of the image transfer for fancni:
```bash
ROCK_FILE="${REPO_ROOT}/fancni_0.1_amd64.rock"
for VM in "${NODE1}" "${NODE2}"; do
  lxc file push "${ROCK_FILE}" "${VM}/tmp/fancni_0.1_amd64.rock"
  lxc_exec "${VM}" "${CTR_BIN}" --address "${CTR_SOCK}" -n "${CTR_NS}" \
    images import --all-platforms /tmp/fancni_0.1_amd64.rock
  lxc_exec "${VM}" "${CTR_BIN}" --address "${CTR_SOCK}" -n "${CTR_NS}" \
    images tag "fancni:0.1" "ghcr.io/ktsakalozos-canonical/fancni:latest"
done
```

### External images (nginx, kube-router) — keep docker save/import as-is
These still come from Docker Hub and are imported the same way.

## Non-goals
- Do not change test logic (phases 8-15)
- Do not change VM creation or k8s setup phases
