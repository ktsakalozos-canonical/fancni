# Task: Update release-latest.yml for rockcraft

## Context
The project now produces a single rock instead of two Docker images. The release workflow needs to build the rock and push it to GHCR using skopeo (since rocks are OCI archives, not Docker images).

## Objective
Replace Docker build+push with rockcraft pack + skopeo copy in `.github/workflows/release-latest.yml`, and capitalize job names.

## Changes

### Job name — capitalize
Add `name: Release` under `jobs.release:`.

### Replace image build and push steps
Remove both existing "Build and push" steps (fancni and fancni-init) and replace with:

1. **Install Rockcraft**:
   ```yaml
   - name: Install Rockcraft
     run: sudo snap install rockcraft --classic
   ```

2. **Build rock**:
   ```yaml
   - name: Build rock
     run: rockcraft pack
   ```

3. **Push rock to GHCR** using skopeo (available via `sudo apt-get install skopeo` or `snap install skopeo`). Actually, rockcraft bundles skopeo as `rockcraft.skopeo`, so no extra install needed:
   ```yaml
   - name: Push rock to GHCR
     env:
       REGISTRY: ghcr.io/${{ github.repository_owner }}
     run: |
       rockcraft.skopeo --insecure-policy copy \
         --dest-creds "${{ github.actor }}:${{ secrets.GITHUB_TOKEN }}" \
         oci-archive:fancni_0.1_amd64.rock \
         "docker://${REGISTRY}/fancni:${SHA}"
       rockcraft.skopeo --insecure-policy copy \
         --dest-creds "${{ github.actor }}:${{ secrets.GITHUB_TOKEN }}" \
         oci-archive:fancni_0.1_amd64.rock \
         "docker://${REGISTRY}/fancni:latest"
   ```

4. **Remove** the `docker/login-action` step — skopeo handles auth via `--dest-creds` directly. Keep the `REGISTRY` and `SHA` env vars at the job level.

### Final step list should be:
1. Checkout (unchanged)
2. Install Rockcraft
3. Build rock
4. Push rock to GHCR

## Non-goals
- Do not change workflow triggers or permissions.
