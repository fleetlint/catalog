#!/usr/bin/env python3
"""Regenerate presets/strict.yaml from minimal and recommended.

strict only overrides: every warning becomes an error, every info a warning,
and the security rules take no exceptions. Run after adding a rule; the
fleetlint test TestStrictAndOssPresets fails when this is stale.
"""
import yaml

NO_EXCEPTIONS = {
    "hooks/secret-scan", "ci/actions-pinned", "ci/least-privilege", "ci/no-pull-request-target",
    "release/signing", "release/provenance", "release/checksums", "release/signed-tags",
    "deps/vulnerability-scan", "repo/lockfile-committed", "repo/no-tracked-env",
    "public/license", "public/no-agent-files", "public/no-ai-attribution",
}
RAISE = {"warning": "error", "info": "warning"}

rules = []
for preset in ("minimal", "recommended"):
    with open(f"presets/{preset}.yaml", encoding="utf-8") as f:
        rules += yaml.safe_load(f)["rules"]
missing = NO_EXCEPTIONS - {r["id"] for r in rules}
if missing:
    raise SystemExit(f"unknown rules in NO_EXCEPTIONS: {sorted(missing)}")

with open("presets/strict.yaml", encoding="utf-8") as f:
    text = f.read()
head = text[: text.index("overrides:\n") + len("overrides:\n")]
lines = []
for r in rules:
    parts = []
    if r["severity"] in RAISE:
        parts.append(f"severity: {RAISE[r['severity']]}")
    if r["id"] in NO_EXCEPTIONS:
        parts.append("exceptions: false")
    if parts:
        lines.append(f"  {r['id']}: {{ {', '.join(parts)} }}\n")
with open("presets/strict.yaml", "w", encoding="utf-8") as f:
    f.write(head + "".join(lines))
print(f"{len(lines)} overrides")
