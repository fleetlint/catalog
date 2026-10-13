#!/usr/bin/env bash
# Pushes this catalog version as one OCI artifact: every file of presets/,
# rules/, stacks/, templates/ and aliases.yaml, titled with its path, the way
# fleetlint reads `catalog: {oci: <ref>}`. Prints the manifest digest.
# Usage: scripts/oci-push.sh <registry>/<repository>:<tag> [oras push flags...]
set -euo pipefail
ref=${1:?usage: scripts/oci-push.sh <registry>/<repository>:<tag> [oras push flags...]}
shift
cd "$(dirname "$0")/.."
files=(aliases.yaml)
while IFS= read -r f; do files+=("$f"); done < <(find presets rules stacks templates -type f | LC_ALL=C sort)
# fleetlint reads at most 512 files from an artifact.
if (( ${#files[@]} > 512 )); then
  echo "oci-push: ${#files[@]} files; fleetlint reads at most 512 from an artifact" >&2
  exit 1
fi
oras push "$@" --artifact-type application/vnd.fleetlint.catalog.source.v1 \
  --annotation "org.opencontainers.image.version=${ref##*:}" \
  --format go-template='{{.digest}}' "$ref" "${files[@]}"
