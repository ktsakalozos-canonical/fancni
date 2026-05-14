# Task: Fix rockcraft.yaml issues

## Context
Code review identified that `permissions` is not a valid rockcraft keyword. Also `go-generate: []` is unnecessary noise.

## Objective
Fix the `rockcraft.yaml` so it would pass `rockcraft pack` validation.

## Changes needed

1. **Remove `permissions` block** from the `scripts` part and replace with `override-build` that ensures files are executable:
   ```yaml
   scripts:
     plugin: dump
     source: deploy/scripts
     organize:
       install-cni.sh: install-cni.sh
       init-node.sh: init-node.sh
     override-build: |
       craftctl default
       chmod 0755 ${CRAFT_PART_INSTALL}/install-cni.sh
       chmod 0755 ${CRAFT_PART_INSTALL}/init-node.sh
     stage:
       - install-cni.sh
       - init-node.sh
   ```

2. **Remove `go-generate: []`** from the `fancni-binary` part (unnecessary, default behavior).

## Non-goals
- Do not change anything else in the file.
