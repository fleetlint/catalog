#!/bin/sh
# Writes the GitHub release notes of a catalog version to stdout:
#   scripts/release-notes.sh <tag> <digest> <section.md>
# section.md is the version's section of CHANGELOG.md, cut by the caller.
# The engine level comes from presets/minimal.yaml, or ENGINE when set.
# ATTESTED=0 leaves out the SBOM and provenance (versions before v0.3.1).
set -eu

tag=$1
digest=$2
section=$3
[ -s "$section" ] || {
  echo "release-notes.sh: no CHANGELOG.md section for $tag" >&2
  exit 1
}
engine=${ENGINE:-$(awk '/^[[:space:]]*engine:/ { gsub(/[^0-9]/, ""); print; exit }' presets/minimal.yaml)}
case $engine in
  3) needs=0.4.0 ;;
  4) needs=0.5.0 ;;
  5) needs=0.6.0 ;;
  6) needs=0.7.0 ;;
  7) needs=0.8.0 ;;
  *)
    echo "release-notes.sh: no fleetlint version known for engine level '$engine'" >&2
    exit 1
    ;;
esac
ref=ghcr.io/fleetlint/catalog@$digest

sed '/./,$!d' "$section"
cat <<NOTES

---

**Engine level $engine**: needs fleetlint $needs or newer.

**Pin this version** in \`.fleetlint.yaml\`, by digest from the registry:

\`\`\`yaml
catalog: { oci: $ref }
\`\`\`

or by tag from git: \`catalog: { version: $tag }\`.

As a Go module: \`go get github.com/fleetlint/catalog@$tag\`.

**Verify** the signature of the published artifact:

\`\`\`sh
cosign verify $ref \\
  --certificate-identity-regexp '^https://github\.com/fleetlint/catalog/\.github/workflows/publish\.yml@refs/tags/v' \\
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
NOTES
if [ "${ATTESTED:-1}" = 1 ]; then
  cat <<NOTES
gh attestation verify oci://$ref --repo fleetlint/catalog   # SLSA build provenance
\`\`\`

The SPDX SBOM of the tree is attested to the same digest (\`cosign verify-attestation --type spdxjson\`).
NOTES
else
  echo '```'
fi
