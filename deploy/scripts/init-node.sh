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
    nsenter --target 1 --net -- chroot /host fanctl up -o "$OVERLAY_NETWORK" -u "${HOST_IP}/${UNDERLAY_PREFIX}"
fi

# Disable ICMP redirects on fan bridge (packets legitimately arrive and leave on same interface)
nsenter --target 1 --net -- sysctl -w "net.ipv4.conf.${BRIDGE_NAME}.send_redirects=0"
# Disable reverse path filtering on fan bridge.
# NOTE: all.rp_filter must also be 0 because Linux uses max(all, per-iface) as effective value.
nsenter --target 1 --net -- sysctl -w "net.ipv4.conf.${BRIDGE_NAME}.rp_filter=0"
nsenter --target 1 --net -- sysctl -w "net.ipv4.conf.all.rp_filter=0"
# Ensure IP forwarding is enabled
nsenter --target 1 --net -- sysctl -w "net.ipv4.ip_forward=1"

# 6. Compute pod subnet for iptables
IFS='.' read -r a b c d <<< "$HOST_IP"
POD_SUBNET="${OVERLAY_FIRST}.${c}.${d}.0/24"

# 7. Set iptables rules (idempotent via -C check)
# Use iptables-legacy to match kube-proxy which also uses iptables-legacy on this platform.
for RULE in "-s $POD_SUBNET -j ACCEPT" "-d $POD_SUBNET -j ACCEPT"; do
    if ! nsenter --target 1 --net -- iptables-legacy -C FORWARD $RULE 2>/dev/null; then
        nsenter --target 1 --net -- iptables-legacy -I FORWARD 1 $RULE
    fi
done

MASQ_RULE="-s $POD_SUBNET ! -o $BRIDGE_NAME -j MASQUERADE"
if ! nsenter --target 1 --net -- iptables-legacy -t nat -C POSTROUTING $MASQ_RULE 2>/dev/null; then
    nsenter --target 1 --net -- iptables-legacy -t nat -A POSTROUTING $MASQ_RULE
fi

echo "Node initialization complete."
