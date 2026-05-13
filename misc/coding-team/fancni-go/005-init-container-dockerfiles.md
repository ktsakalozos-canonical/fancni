# Task 005: Init Container & Dockerfiles

## Context

The CNI binary is complete. This task creates the container images and init scripts needed to deploy fancni on Kubernetes nodes via a DaemonSet.

## Objective

Create two Dockerfiles and the init script that together:
1. Build the CNI binary into a minimal image
2. Prepare each node: install ubuntu-fan, load kernel modules, create fan bridge, set iptables rules
3. Install the CNI binary + config onto the host filesystem

## Scope

### Directory structure

```
deploy/
  docker/
    Dockerfile.cni           # Multi-stage: build Go binary, copy to install image
    Dockerfile.init          # Ubuntu-based init container
  scripts/
    init-node.sh             # Node initialization script
    install-cni.sh           # CNI binary + config installation script
```

### deploy/docker/Dockerfile.cni

Multi-stage build:

**Stage 1: builder**
- Base: `golang:1.22-bookworm`
- Copy source, run `go build -o /fancni ./cmd/fancni/`

**Stage 2: install**
- Base: `busybox:1.36` (minimal, just needs `cp` and `sh`)
- Copy the built binary from stage 1
- Copy `deploy/scripts/install-cni.sh`
- Entrypoint: `install-cni.sh`

### deploy/docker/Dockerfile.init

- Base: `ubuntu:24.04`
- Install: `iptables`, `iproute2`, `kmod` (for modprobe)
- Copy `deploy/scripts/init-node.sh`
- Entrypoint: `init-node.sh`

### deploy/scripts/init-node.sh

This runs as a privileged init container with these host mounts:
- `/` → `/host` (host root filesystem, for chroot)
- `/run` → `/host/run` (for systemd/dbus if needed by apt)

The script receives configuration via environment variables (set by the Helm chart):
- `OVERLAY_NETWORK` (default: `240.0.0.0/8`)
- `UNDERLAY_PREFIX` (default: `16`)

**Script logic:**

```bash
#!/bin/bash
set -e

OVERLAY_NETWORK="${OVERLAY_NETWORK:-240.0.0.0/8}"
UNDERLAY_PREFIX="${UNDERLAY_PREFIX:-16}"

# 1. Install ubuntu-fan on host if not present
if ! chroot /host which fanctl >/dev/null 2>&1; then
    echo "Installing ubuntu-fan on host..."
    chroot /host apt-get update -qq
    chroot /host apt-get install -y -qq ubuntu-fan
fi

# 2. Load kernel modules on host
nsenter --target 1 --mount -- modprobe vxlan || true
nsenter --target 1 --mount -- modprobe ipip || true

# 3. Detect host IP (using the host's network namespace)
HOST_IP=$(nsenter --target 1 --net -- ip route get 1.1.1.1 | awk '{for(i=1;i<=NF;i++) if($i=="src") print $(i+1); exit}')
if [ -z "$HOST_IP" ]; then
    echo "ERROR: Could not detect host IP"
    exit 1
fi
echo "Host IP: $HOST_IP"

# 4. Compute fan bridge name from overlay
OVERLAY_FIRST=$(echo "$OVERLAY_NETWORK" | cut -d. -f1)
BRIDGE_NAME="fan-${OVERLAY_FIRST}"

# 5. Create fan bridge if not present
if ! nsenter --target 1 --net -- ip link show "$BRIDGE_NAME" >/dev/null 2>&1; then
    echo "Creating fan bridge..."
    chroot /host fanctl up -o "$OVERLAY_NETWORK" -u "${HOST_IP}/${UNDERLAY_PREFIX}"
fi

# 6. Compute pod subnet for iptables
IFS='.' read -r a b c d <<< "$HOST_IP"
POD_SUBNET="${OVERLAY_FIRST}.${c}.${d}.0/24"

# 7. Set iptables rules (idempotent via -C check)
for RULE in "-s $POD_SUBNET -j ACCEPT" "-d $POD_SUBNET -j ACCEPT"; do
    if ! iptables -C FORWARD $RULE 2>/dev/null; then
        iptables -A FORWARD $RULE
    fi
done

MASQ_RULE="-s $POD_SUBNET ! -o $BRIDGE_NAME -j MASQUERADE"
if ! iptables -t nat -C POSTROUTING $MASQ_RULE 2>/dev/null; then
    iptables -t nat -A POSTROUTING $MASQ_RULE
fi

echo "Node initialization complete."
```

### deploy/scripts/install-cni.sh

This runs as another init container. Simpler — just copies files.

Host mounts:
- `/opt/cni/bin` → `/host/opt/cni/bin`
- `/etc/cni/net.d` → `/host/etc/cni/net.d`

Configuration via env vars:
- `CNI_CONF` — the full JSON CNI config (injected from ConfigMap by Helm)

```bash
#!/bin/sh
set -e

# Copy CNI binary
echo "Installing fancni binary..."
cp /fancni /host/opt/cni/bin/fancni
chmod +x /host/opt/cni/bin/fancni

# Write CNI config
echo "Writing CNI config..."
echo "$CNI_CONF" > /host/etc/cni/net.d/10-fancni.conf

echo "CNI installation complete."

# Sleep forever to keep the pod alive (the DaemonSet pod needs to stay running).
# This container's job is done, but the pod must remain so kubelet knows the CNI is installed.
exec sleep infinity
```

Note: `sleep infinity` in the install container keeps the DaemonSet pod running. An alternative is a separate pause container, but using the installer itself is simpler.

### Makefile updates

Add these targets to the existing Makefile:

```makefile
docker-build-cni:
	docker build -t fancni:latest -f deploy/docker/Dockerfile.cni .

docker-build-init:
	docker build -t fancni-init:latest -f deploy/docker/Dockerfile.init .

docker-build: docker-build-cni docker-build-init
```

## Non-goals

- No Helm chart (Phase 6)
- No CI/CD pipeline
- No image push to registry

## Constraints

- Init container must be privileged (needs host namespace access, chroot, iptables)
- Scripts must be idempotent (safe to re-run on node restart)
- Use `set -e` in all scripts for fail-fast behavior
- The install container must keep running (`sleep infinity`) so the DaemonSet pod stays alive
