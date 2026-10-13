# Changelog

All notable changes to this project are documented here.
Format: Keep a Changelog. Versioning: Semantic Versioning.

## [0.1.0] - 2026-10-13

The first release of the catalog on its own. fleetlint 0.3 and older embedded the rules under the ids listed as renamed below.

### Added

- The presets `minimal`, `recommended`, `slop`, `strict` and `oss`, the rule library (94 rules and 12 capabilities with 129 per-stack implementations), the stack definitions and the fix templates, moved here from the fleetlint repository. fleetlint embeds this module.
- **Capabilities.** A requirement that every stack meets with its own tools is one capability with an implementation per stack at `rules/<family>/<name>/<stack>.yaml`; `any.yaml` is the stack-neutral implementation, and an implementation with `when: "false"` says the requirement does not apply to its stack. Capabilities: `lint/configured`, `lint/strict`, `types/strict`, `format/configured`, `repo/lockfile-committed`, `repo/toolchain-pinned`, `deps/vulnerability-scan`, `deps/license-check`, `deps/locked-install`, `deps/unused-check`, `quality/coverage-threshold`, `release/reproducible-build`. Selecting, overriding or excepting a capability reaches all its implementations.
- **Stacks as data.** `stacks/<id>.yaml` for cpp, dotnet, flutter, go, hugo, java, kotlin, node, php, python, ruby and rust: detection markers, dev container feature or image, the image targets run in, Renovate lockfile maintenance.
- New rules: `ci/no-script-injection`, `ci/no-pipe-to-shell`, `ci/checkout-no-persist-credentials`, `ci/workflow-lint`, `ci/sast`, `hooks/revs-pinned`, `release/trusted-publishing`, `release/hosted-builder`, `deps/update-coverage`, `deps/update-cooldown`, `slop/forced-unwraps`; new capabilities `format/configured`, `types/strict`, `deps/locked-install`, `deps/unused-check`, `release/reproducible-build`.
- GitLab CI: every rule that reads CI understands `.gitlab-ci.yml` (jobs, scripts, `timeout:`, `interruptible`, includes pinned to a commit, release pipelines on `$CI_COMMIT_TAG`); Forgejo Actions workflows are read too. lefthook is understood by `hooks/secret-scan`, `hooks/conventional-commits` and `hooks/pre-push-check`.
- Every tag is published as a signed OCI artifact, `ghcr.io/fleetlint/catalog:<tag>`, holding the whole catalog version; a repository pins it with `catalog: {oci: ghcr.io/fleetlint/catalog:<tag>}`. The check workflow builds and loads the artifact on every change.
- `aliases.yaml` maps every renamed rule id to the ids that replace it; fleetlint reads an old id in a configuration as those, with a warning.
- `deps/update-automation` has a fix: `fleetlint fix` writes a renovate.json with only the settings the repository's stacks need.
- Templates for a justfile and a Taskfile next to the Makefile, and workflow fragments that install those runners.

### Changed

- **Renamed.** `lint/<stack>-config` is `lint/configured/<stack>`, `lint/hugo-strict-build` is `lint/configured/hugo`; `lint/go-linters-enabled`, `lint/python-rules-enabled`, `lint/rust-lints-denied` and `lint/kotlin-detekt-strict` are `lint/strict/<stack>`; `lint/node-strict-flags` and `lint/flutter-strict-modes` are `types/strict/<stack>`, and the tsconfig and mypy parts of `lint/node-config` and `lint/python-config` moved there; `quality/no-slop-markers` is `slop/marketing-words` (info, in `slop`).
- **Presets.** `public/no-agent-files` and `public/no-ai-attribution` moved from `minimal` to `slop`. `minimal` adds `slop/conflict-markers`, `repo/no-large-files` and `public/license`. `recommended` adds `format/*` and `types/*` and lists its public rules. `strict` also forbids exceptions on `ci/no-script-injection`, `ci/no-pipe-to-shell`, `hooks/revs-pinned` and `release/trusted-publishing`. Presets state `engine: 3` and no version of their own: they take the release tag.
- The lint, format and type capabilities apply at tiers 1 and 2; the lint rules before applied at every tier. Lint rules find a configuration in a parent directory, as the tools do, so workspace members inherit the root's.
- `public/no-ai-attribution` matches AI tool and bot addresses and generator footers, never a person's name. `public/no-agent-files` also covers `GEMINI.md`, `.github/copilot-instructions.md`, `.cursorrules`, `.windsurfrules`, `.clinerules`, `.aider*` and `.cursor/`.
- Vulnerability scanners, license checks and coverage gates count only where they are run: an `echo`, a comment or an install line no longer satisfies those rules, and the templates' TODO lines name no tool.
- `hooks/secret-scan` needs a secret scanner (gitleaks, trufflehog, ggshield, detect-secrets, secretlint); `detect-private-key` alone is a partial. `ci/check-workflow` and `hooks/pre-push-check` no longer accept `check-fast`. `ci/actions-pinned` needs a digest for `docker://` and a commit for GitLab includes. `repo/toolchain-pinned` needs the tool named with a version (and a JDK pin for Java and Kotlin); `repo/lockfile-committed` no longer counts an unhashed requirements.txt.
- `release/sbom` is a warning and does not apply when a release only publishes to a package registry or when git-cliff is the only release tool.
- `repo/no-tracked-env` and `repo/no-tracked-junk` run once at the root. The per-stack capabilities run in every workspace member and nested project, and accept what the repository provides above it: a workspace's root lockfile, a pin in a parent directory, a scanner or coverage gate in the root's CI or task runner.
- A Gradle build is `java` when it applies the Java plugins and neither Kotlin nor Android; every other Gradle build (Kotlin, Android, Kotlin Multiplatform) is `kotlin`.
- `ci/least-privilege` checks GitHub Actions workflows only: Gitea and Forgejo do not take the token's scope from the workflow file.
- `repo/readme-present` accepts a heading on running, developing or building the project too.
- `repo/no-tracked-env` and `repo/no-tracked-junk` pass `.env` files, logs and dependency directories under test and fixture directories: they are test inputs.
- A Hugo site's `go.mod` (Hugo modules) does not make it a Go project; Go sources next to it do.
- Slop heuristics report fewer false positives: `model.eval()`, `iter().skip()`, "Enter your password", mdBook's `SUMMARY.md`, pytest's `test_*.py`, `isNew`, versioned API files, Git's `Revert "..."` subjects and the initial commit no longer count; conflict markers are found in every file type.
- Templates: ESLint flat config with `defineConfig`, PHPStan 2 and PHP CS Fixer rule sets, clang-tidy header filter, ruff without a hardcoded target version, current Dart analyzer options, `.editorconfig` with tabs for Go and two spaces for Dart, `ci` and `style` commit types, npm publishing through trusted publishing.

### Fixed

- `ci/concurrency-cancel` no longer errors on `concurrency` without `cancel-in-progress`. `hooks/conventional-commits` no longer passes on `no-commit-to-branch`. `containers/non-root-user` reads the final stage only.
- `slop/focused-tests` reads test files only, and looks for `fdescribe`, `fit` and `.only` in JavaScript and TypeScript and for `solo: true` in Dart; a Go function called `fit` was an error before.
- `repo/no-tracked-env` accepts encrypted environment files (`.sops`, `.enc`, `.age`, `.gpg`, `.asc`).
- The agent instructions of `ci/least-privilege` and `ci/concurrency-cancel` showed an escaped `\n`; they now give the YAML in flow form.
