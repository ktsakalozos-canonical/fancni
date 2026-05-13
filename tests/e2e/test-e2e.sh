#!/usr/bin/env bash
# E2E test for fancni CNI plugin using LXC VMs and Canonical K8s.
set -euo pipefail

# ---------------------------------------------------------------------------
# Config
# ---------------------------------------------------------------------------
NODE1="fancni-e2e-node-1"
NODE2="fancni-e2e-node-2"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
CTR_BIN="/snap/k8s/current/bin/ctr"
CTR_SOCK="/run/containerd/containerd.sock"
CTR_NS="k8s.io"

# ---------------------------------------------------------------------------
# Cleanup / flags
# ---------------------------------------------------------------------------
NO_CLEANUP="${FANCNI_E2E_NO_CLEANUP:-0}"
if [[ "${1:-}" == "--no-cleanup" ]]; then
  NO_CLEANUP=1
fi

cleanup() {
  if [[ "$NO_CLEANUP" == "1" ]]; then
    echo "[e2e] --no-cleanup set; skipping VM deletion (${NODE1}, ${NODE2})"
    return
  fi
  echo "[e2e] Cleaning up VMs..."
  lxc delete --force "${NODE1}" "${NODE2}" 2>/dev/null || true
}
trap cleanup EXIT

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------
log() { echo "[e2e] $*"; }

# wait_for <timeout_seconds> <description> <cmd…>
wait_for() {
  local timeout="$1"; shift
  local desc="$1"; shift
  local elapsed=0
  local interval=5
  log "Waiting for: ${desc} (timeout ${timeout}s)"
  while ! "$@" &>/dev/null; do
    if (( elapsed >= timeout )); then
      echo "[e2e] TIMEOUT waiting for: ${desc}" >&2
      return 1
    fi
    sleep "${interval}"
    elapsed=$(( elapsed + interval ))
  done
  log "Ready: ${desc}"
}

# retry <max_attempts> <delay_seconds> <cmd…>
retry() {
  local max="$1"; shift
  local delay="$1"; shift
  local attempt=1
  while true; do
    if "$@"; then
      return 0
    fi
    if (( attempt >= max )); then
      log "Command failed after ${max} attempts: $*"
      return 1
    fi
    log "  Attempt ${attempt}/${max} failed, retrying in ${delay}s..."
    sleep "${delay}"
    (( attempt++ ))
  done
}

lxc_exec() {
  local vm="$1"; shift
  lxc exec "${vm}" -- "$@"
}

# ---------------------------------------------------------------------------
# Phase 1 – Create VMs
# ---------------------------------------------------------------------------
log "=== Phase 1: Creating LXC VMs ==="
lxc delete --force "${NODE1}" "${NODE2}" 2>/dev/null || true

lxc launch ubuntu:24.04 "${NODE1}" --vm \
  --config limits.cpu=2 \
  --config limits.memory=4GiB

lxc launch ubuntu:24.04 "${NODE2}" --vm \
  --config limits.cpu=2 \
  --config limits.memory=4GiB

# ---------------------------------------------------------------------------
# Phase 2 – Wait for cloud-init
# ---------------------------------------------------------------------------
log "=== Phase 2: Waiting for cloud-init ==="
wait_for 300 "cloud-init ${NODE1}" \
  lxc exec "${NODE1}" -- cloud-init status --wait
wait_for 300 "cloud-init ${NODE2}" \
  lxc exec "${NODE2}" -- cloud-init status --wait

# ---------------------------------------------------------------------------
# Phase 3 – Install k8s snap
# ---------------------------------------------------------------------------
log "=== Phase 3: Installing k8s snap ==="
retry 5 30 lxc_exec "${NODE1}" snap install k8s --classic --channel=1.35-classic/stable
retry 5 30 lxc_exec "${NODE2}" snap install k8s --classic --channel=1.35-classic/stable

# ---------------------------------------------------------------------------
# Phase 4 – Bootstrap node-1 without CNI
# ---------------------------------------------------------------------------
log "=== Phase 4: Bootstrapping ${NODE1} (no CNI) ==="
config="/tmp/fancni-bootstrap-config.yaml"
cat <<EOF > "${config}"
cluster-config:
  network:
    enabled: false
EOF
lxc file push "${config}" "${NODE1}/tmp/bootstrap-config.yaml"
rm -f "${config}"
lxc_exec "${NODE1}" sudo k8s bootstrap --file /tmp/bootstrap-config.yaml

wait_for 120 "k8s API responding on ${NODE1}" \
  lxc exec "${NODE1}" -- sudo k8s kubectl get nodes

# ---------------------------------------------------------------------------
# Phase 5 – Join node-2
# ---------------------------------------------------------------------------
log "=== Phase 5: Joining ${NODE2} to the cluster ==="
JOIN_TOKEN="$(lxc_exec "${NODE1}" sudo k8s get-join-token "${NODE2}" --worker)"
log "Join token obtained"
lxc_exec "${NODE2}" sudo k8s join-cluster "${JOIN_TOKEN}"

wait_for 300 "both nodes visible" bash -c \
  "lxc exec ${NODE1} -- sudo k8s kubectl get nodes 2>/dev/null | grep -c Ready | grep -qE '^[2-9]'"

# ---------------------------------------------------------------------------
# Phase 6 – Build fancni images on the host
# ---------------------------------------------------------------------------
log "=== Phase 6: Building fancni images on host ==="
make -C "${REPO_ROOT}" docker-build

# ---------------------------------------------------------------------------
# Phase 7 – Transfer images into VMs
# ---------------------------------------------------------------------------
log "=== Phase 7: Transferring images into VMs ==="
docker pull nginx:latest
docker pull cloudnativelabs/kube-router:v2.9.0
for IMAGE in "fancni:latest" "fancni-init:latest" "nginx:latest" "cloudnativelabs/kube-router:v2.9.0"; do
  for VM in "${NODE1}" "${NODE2}"; do
    log "  Importing ${IMAGE} into ${VM}..."
    docker save "${IMAGE}" | lxc exec "${VM}" -- \
      "${CTR_BIN}" \
        --address "${CTR_SOCK}" \
        -n "${CTR_NS}" \
        images import -
  done
done

# Re-tag fancni images to match Helm chart references
for VM in "${NODE1}" "${NODE2}"; do
  log "  Re-tagging images in ${VM}..."
  lxc exec "${VM}" -- "${CTR_BIN}" --address "${CTR_SOCK}" -n "${CTR_NS}" \
    images tag "docker.io/library/fancni:latest" "ghcr.io/ktsakalozos-canonical/fancni:latest"
  lxc exec "${VM}" -- "${CTR_BIN}" --address "${CTR_SOCK}" -n "${CTR_NS}" \
    images tag "docker.io/library/fancni-init:latest" "ghcr.io/ktsakalozos-canonical/fancni-init:latest"
done

# ---------------------------------------------------------------------------
# Phase 8 – Install fancni via Helm
# ---------------------------------------------------------------------------
log "=== Phase 8: Installing fancni via Helm ==="
lxc_exec "${NODE1}" mkdir -p /tmp/fancni-chart
lxc file push --recursive "${REPO_ROOT}/deploy/helm/fancni/" "${NODE1}/tmp/fancni-chart/"
lxc_exec "${NODE1}" sudo k8s helm install fancni /tmp/fancni-chart/fancni --namespace kube-system

log "Waiting for fancni DaemonSet pods to be Running..."
wait_for 300 "fancni DaemonSet ready" bash -c \
  "lxc exec ${NODE1} -- sudo k8s kubectl -n kube-system get pods -l app=fancni 2>/dev/null \
   | awk 'NR>1{print \$3}' | grep -v Running | wc -l | grep -q '^0$'"

# ---------------------------------------------------------------------------
# Phase 9 – Deploy nginx
# ---------------------------------------------------------------------------
log "=== Phase 9: Deploying nginx (4 replicas) ==="
lxc_exec "${NODE1}" sudo k8s kubectl create deployment nginx-e2e \
  --image=nginx:latest --replicas=4

wait_for 300 "4 nginx pods Running" bash -c \
  "lxc exec ${NODE1} -- sudo k8s kubectl get pods -l app=nginx-e2e 2>/dev/null \
   | awk 'NR>1{print \$3}' | grep -c Running | grep -q '^4$'"

# ---------------------------------------------------------------------------
# Phase 10 – Assert pods got fancni IPs (240.x.x.x)
# ---------------------------------------------------------------------------
log "=== Phase 10: Asserting pod IPs are in 240.0.0.0/8 ==="
POD_IPS=$(lxc_exec "${NODE1}" sudo k8s kubectl get pods -l app=nginx-e2e \
  -o jsonpath='{.items[*].status.podIP}')

PASS_COUNT=0
FAIL_COUNT=0
FAILED_IPS=()

for IP in ${POD_IPS}; do
  if [[ "${IP}" == 240.* ]]; then
    log "  PASS: pod IP ${IP} is in 240.0.0.0/8"
    (( PASS_COUNT++ )) || true
  else
    log "  FAIL: pod IP ${IP} is NOT in 240.0.0.0/8"
    (( FAIL_COUNT++ )) || true
    FAILED_IPS+=("${IP}")
  fi
done

# ---------------------------------------------------------------------------
# Phase 11 – Assert HTTP 200 from each pod IP
# ---------------------------------------------------------------------------
log "=== Phase 11: Asserting HTTP 200 from pod IPs ==="
HTTP_PASS=0
HTTP_FAIL=0
FAILED_HTTP_IPS=()

for IP in ${POD_IPS}; do
  HTTP_CODE="000"
  for attempt in $(seq 1 5); do
    HTTP_CODE=$(lxc_exec "${NODE1}" \
      curl -s -o /dev/null -w '%{http_code}' --max-time 5 "http://${IP}" 2>/dev/null) || true
    if [[ "${HTTP_CODE}" == "200" ]]; then
      break
    fi
    log "  Attempt ${attempt}/5 for ${IP} returned ${HTTP_CODE}, retrying in 10s..."
    sleep 10
  done
  if [[ "${HTTP_CODE}" == "200" ]]; then
    log "  PASS: http://${IP} -> ${HTTP_CODE}"
    (( HTTP_PASS++ )) || true
  else
    log "  FAIL: http://${IP} -> ${HTTP_CODE} (expected 200)"
    (( HTTP_FAIL++ )) || true
    FAILED_HTTP_IPS+=("${IP}")
  fi
done

# ---------------------------------------------------------------------------
# Phase 12 – Verify kube-router pods running
# ---------------------------------------------------------------------------
log "=== Phase 12: Verifying kube-router pods are Running ==="
wait_for 300 "kube-router pods Running on both nodes" bash -c \
  "lxc exec ${NODE1} -- sudo k8s kubectl -n kube-system get pods -l app=kube-router 2>/dev/null \
   | awk 'NR>1{print \$3}' | grep -v Running | wc -l | grep -q '^0$'"

# ---------------------------------------------------------------------------
# Phase 13 – Test deny-all NetworkPolicy blocks traffic
# ---------------------------------------------------------------------------
log "=== Phase 13: Testing deny-all NetworkPolicy blocks traffic ==="

# Create namespace and workloads
lxc_exec "${NODE1}" sudo k8s kubectl create namespace netpol-test || true

lxc_exec "${NODE1}" sudo k8s kubectl -n netpol-test run web \
  --image=nginx:latest --labels="app=web" --restart=Never

lxc_exec "${NODE1}" sudo k8s kubectl -n netpol-test run client \
  --image=nginx:latest --labels="app=client" --restart=Never

wait_for 120 "web pod Running in netpol-test" bash -c \
  "lxc exec ${NODE1} -- sudo k8s kubectl -n netpol-test get pod web 2>/dev/null \
   | awk 'NR>1{print \$3}' | grep -q Running"

wait_for 120 "client pod Running in netpol-test" bash -c \
  "lxc exec ${NODE1} -- sudo k8s kubectl -n netpol-test get pod client 2>/dev/null \
   | awk 'NR>1{print \$3}' | grep -q Running"

WEB_POD_IP=$(lxc_exec "${NODE1}" sudo k8s kubectl -n netpol-test get pod web \
  -o jsonpath='{.status.podIP}')
log "web pod IP: ${WEB_POD_IP}"

# Verify connectivity before policy
NETPOL_PASS=1
PRE_CODE="000"
for attempt in $(seq 1 5); do
  PRE_CODE=$(lxc_exec "${NODE1}" \
    sudo k8s kubectl -n netpol-test exec client -- \
    curl -s -o /dev/null -w '%{http_code}' --max-time 5 "http://${WEB_POD_IP}" 2>/dev/null) || true
  if [[ "${PRE_CODE}" == "200" ]]; then
    break
  fi
  log "  Pre-policy attempt ${attempt}/5 returned ${PRE_CODE}, retrying in 5s..."
  sleep 5
done

if [[ "${PRE_CODE}" == "200" ]]; then
  log "  PASS: pre-policy curl returned 200"
else
  log "  FAIL: pre-policy curl returned ${PRE_CODE} (expected 200)"
  NETPOL_PASS=0
fi

# Apply deny-all ingress policy
lxc_exec "${NODE1}" sudo k8s kubectl apply -f - <<'EOF'
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all-ingress
  namespace: netpol-test
spec:
  podSelector:
    matchLabels:
      app: web
  policyTypes:
  - Ingress
EOF

log "Waiting 15s for kube-router to sync deny-all rules..."
sleep 15

# Verify connectivity is now blocked (retry a few times to confirm it stays blocked)
DENY_CODE="200"
for attempt in $(seq 1 5); do
  DENY_CODE=$(lxc_exec "${NODE1}" \
    sudo k8s kubectl -n netpol-test exec client -- \
    curl -s -o /dev/null -w '%{http_code}' --max-time 5 "http://${WEB_POD_IP}" 2>/dev/null) || true
  if [[ "${DENY_CODE}" != "200" ]]; then
    break
  fi
  log "  Deny attempt ${attempt}/5 still got ${DENY_CODE}, waiting 10s..."
  sleep 10
done

if [[ "${DENY_CODE}" != "200" ]]; then
  log "  PASS: deny-all policy blocked traffic (got ${DENY_CODE})"
else
  log "  FAIL: deny-all policy did NOT block traffic (still got 200)"
  NETPOL_PASS=0
fi

# ---------------------------------------------------------------------------
# Phase 14 – Test allow NetworkPolicy restores traffic
# ---------------------------------------------------------------------------
log "=== Phase 14: Testing allow NetworkPolicy restores traffic ==="

lxc_exec "${NODE1}" sudo k8s kubectl apply -f - <<'EOF'
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-client
  namespace: netpol-test
spec:
  podSelector:
    matchLabels:
      app: web
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: client
EOF

log "Waiting 15s for kube-router to sync allow rules..."
sleep 15

# Verify connectivity is restored
ALLOW_CODE="000"
for attempt in $(seq 1 5); do
  ALLOW_CODE=$(lxc_exec "${NODE1}" \
    sudo k8s kubectl -n netpol-test exec client -- \
    curl -s -o /dev/null -w '%{http_code}' --max-time 5 "http://${WEB_POD_IP}" 2>/dev/null) || true
  if [[ "${ALLOW_CODE}" == "200" ]]; then
    break
  fi
  log "  Allow attempt ${attempt}/5 returned ${ALLOW_CODE}, retrying in 10s..."
  sleep 10
done

if [[ "${ALLOW_CODE}" == "200" ]]; then
  log "  PASS: allow policy restored traffic (got 200)"
else
  log "  FAIL: allow policy did NOT restore traffic (got ${ALLOW_CODE})"
  NETPOL_PASS=0
fi

# ---------------------------------------------------------------------------
# Phase 15 – Summary (renumbered from 12)
# ---------------------------------------------------------------------------
log "=== Phase 15: Summary ==="
log "IP assertion:          PASS=${PASS_COUNT}     FAIL=${FAIL_COUNT}"
log "HTTP assertion:        PASS=${HTTP_PASS}      FAIL=${HTTP_FAIL}"
log "NetworkPolicy assert:  PASS=$(( NETPOL_PASS == 1 ? 1 : 0 ))  FAIL=$(( NETPOL_PASS == 0 ? 1 : 0 ))"

OVERALL_PASS=1
if (( FAIL_COUNT > 0 )); then
  log "FAILED IPs (not in 240.0.0.0/8): ${FAILED_IPS[*]}"
  OVERALL_PASS=0
fi
if (( HTTP_FAIL > 0 )); then
  log "FAILED HTTP IPs: ${FAILED_HTTP_IPS[*]}"
  OVERALL_PASS=0
fi
if (( NETPOL_PASS == 0 )); then
  log "FAILED: NetworkPolicy enforcement check(s) failed"
  OVERALL_PASS=0
fi

if (( OVERALL_PASS == 1 )); then
  log "=== E2E RESULT: PASS ==="
  exit 0
else
  log "=== E2E RESULT: FAIL ==="
  exit 1
fi
