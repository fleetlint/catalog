#!/usr/bin/env bash
# Loads a pushed catalog version with fleetlint the way a repository pins it
# (`catalog: {oci: <ref>}`) and fails unless every preset resolves from it.
# Usage: scripts/oci-check.sh <registry>/<repository>[:<tag>][@sha256:<digest>]
set -euo pipefail
ref=${1:?usage: scripts/oci-check.sh <ref>}
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
git -C "$work" init -q
printf 'module example.com/x\n\ngo 1.24\n' > "$work/go.mod"
for preset in $(cd "$(dirname "$0")/../presets" && ls ./*.yaml | sed 's|^\./||; s|\.yaml$||'); do
  printf 'version: 1\ncatalog: {oci: "%s"}\nextends: [fleetlint:%s]\n' "$ref" "$preset" > "$work/.fleetlint.yaml"
  out=$(cd "$work" && fleetlint catalog version)
  grep -q "pinned" <<<"$out" || { echo "oci-check: $preset: the run does not use the pinned artifact:" >&2; echo "$out" >&2; exit 1; }
  # Exit 0 or 1 is a run; 2 means the configuration or the catalog did not load.
  status=0
  (cd "$work" && fleetlint check --format json > "$work/$preset.json") || status=$?
  if (( status > 1 )); then
    echo "oci-check: $preset: fleetlint check exited $status" >&2
    exit 1
  fi
  echo "oci-check: $preset loads from $ref"
done
