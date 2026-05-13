# Task: Create docs/configuration.md

## Context
The README has a Configuration section with a Helm values table. We're moving it to its own file.

## Objective
Create `docs/configuration.md` containing the Helm chart configuration reference.

## Content to include
Move the following table from the README (lines 59-74) into this doc:

| Key | Default | Description |
|-----|---------|-------------|
| `overlay.network` | `240.0.0.0/8` | Fan overlay address space |
| `underlay.prefix` | `16` | Underlay subnet prefix length used to derive per-node pod CIDRs |
| `images.cni.repository` | `fancni` | CNI container image name |
| `images.cni.tag` | `latest` | CNI container image tag |
| `images.init.repository` | `fancni-init` | Init container image name |
| `images.init.tag` | `latest` | Init container image tag |
| `cni.name` | `fancni` | CNI plugin name |
| `cni.confFileName` | `10-fancni.conf` | CNI config file name written to `/etc/cni/net.d/` |
| `ipam.dataDir` | `/var/lib/cni/fancni` | Directory for IPAM state files |

Add a brief intro noting these values live in `deploy/helm/fancni/values.yaml`.

## Non-goals
- No usage examples beyond the table
- Don't modify the README yet (task 003 handles that)
