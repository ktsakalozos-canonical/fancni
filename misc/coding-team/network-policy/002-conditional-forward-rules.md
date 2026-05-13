# Task: Make FORWARD ACCEPT iptables rules conditional on network policy

## Context

`deploy/scripts/init-node.sh` inserts blanket `FORWARD -s POD_SUBNET -j ACCEPT` and `FORWARD -d POD_SUBNET -j ACCEPT` rules at position 1 in the FORWARD chain. When kube-router manages network policies, it needs to control the FORWARD chain — these blanket ACCEPTs would short-circuit kube-router's DROP rules.

The MASQUERADE rule in nat/POSTROUTING is fine and must remain regardless (kube-router doesn't manage NAT in firewall-only mode).

## Objective

Make the FORWARD ACCEPT rules conditional: skip them when network policy is enabled (kube-router will manage FORWARD).

## Scope

1. **`deploy/scripts/init-node.sh`** — check an env var `NETWORK_POLICY_ENABLED`. If set to `"true"`, skip the FORWARD ACCEPT rules (lines 51-55). Keep the MASQUERADE rule unconditionally.

2. **`deploy/helm/fancni/templates/daemonset.yaml`** — pass the env var `NETWORK_POLICY_ENABLED` to the `node-init` container, derived from `.Values.networkPolicy.enabled`.

## Implementation detail

In `init-node.sh`, wrap the FORWARD rules loop:
```bash
if [ "${NETWORK_POLICY_ENABLED}" != "true" ]; then
    for RULE in "-s $POD_SUBNET -j ACCEPT" "-d $POD_SUBNET -j ACCEPT"; do
        if ! nsenter --target 1 --net -- iptables-legacy -C FORWARD $RULE 2>/dev/null; then
            nsenter --target 1 --net -- iptables-legacy -I FORWARD 1 $RULE
        fi
    done
fi
```

In the daemonset template, add to the `node-init` container's env:
```yaml
- name: NETWORK_POLICY_ENABLED
  value: {{ .Values.networkPolicy.enabled | quote }}
```

## Non-goals

- Do not modify Go source code
- Do not change the MASQUERADE rule behavior
- Do not remove the dead `internal/iptables` Go package (unrelated cleanup)
