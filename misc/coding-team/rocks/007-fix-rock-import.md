# Task: Fix rock import in e2e test

## Context
`ctr images import` of an OCI archive with a partial ref (just `0.1`) generates a name like `import-2026-05-14:0.1` unless `--base-name` is specified. The e2e test was trying to tag `fancni:0.1` which doesn't exist.

## Objective
Fix the rock import command in the e2e test to use `--base-name` so the image gets a proper name.

## Changes to `tests/e2e/test-e2e.sh`

In Phase 7, change the import + tag commands from:
```bash
lxc_exec "${VM}" "${CTR_BIN}" --address "${CTR_SOCK}" -n "${CTR_NS}" \
  images import --all-platforms /tmp/fancni_0.1_amd64.rock
lxc_exec "${VM}" "${CTR_BIN}" --address "${CTR_SOCK}" -n "${CTR_NS}" \
  images tag "fancni:0.1" "ghcr.io/ktsakalozos-canonical/fancni:latest"
```

To:
```bash
lxc_exec "${VM}" "${CTR_BIN}" --address "${CTR_SOCK}" -n "${CTR_NS}" \
  images import --all-platforms --base-name "ghcr.io/ktsakalozos-canonical/fancni" /tmp/fancni_0.1_amd64.rock
lxc_exec "${VM}" "${CTR_BIN}" --address "${CTR_SOCK}" -n "${CTR_NS}" \
  images tag "ghcr.io/ktsakalozos-canonical/fancni:0.1" "ghcr.io/ktsakalozos-canonical/fancni:latest"
```

## Non-goals
- Do not change anything else in the test.
