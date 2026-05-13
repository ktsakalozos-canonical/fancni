#!/bin/sh
set -e

CNI_CONF_FILE="${CNI_CONF_FILE:-10-fancni.conf}"

# Copy CNI binary
echo "Installing fancni binary..."
cp /fancni /host/opt/cni/bin/fancni
chmod +x /host/opt/cni/bin/fancni

# Write CNI config
echo "Writing CNI config..."
printf '%s' "$CNI_CONF" > "/host/etc/cni/net.d/${CNI_CONF_FILE}"

echo "CNI installation complete."

# Sleep forever to keep the pod alive (the DaemonSet pod needs to stay running).
# This container's job is done, but the pod must remain so kubelet knows the CNI is installed.
exec sleep infinity
