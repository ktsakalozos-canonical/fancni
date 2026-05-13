# Completion Report: Fix cross-node pod connectivity

## Summary

- Changed `iptables -A FORWARD` to `iptables -I FORWARD 1` for both ACCEPT rules so they are inserted at the top of the FORWARD chain, evaluated before the `KUBE-FORWARD` chain's DROP rules.
- Added sysctl settings after the fan bridge creation block: disables ICMP redirects and reverse path filtering on the fan bridge, and explicitly enables IP forwarding.
- The `-C` idempotency checks remain unchanged (they check rule existence regardless of position).

## Files Changed

- `deploy/scripts/init-node.sh`

## Notable Tradeoffs or Risks

- None. Changes are minimal and targeted. The sysctl settings run unconditionally on every init (not just on bridge creation), which is safe and idempotent.
