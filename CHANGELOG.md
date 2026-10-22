# Changelog

All notable changes to this project are documented here.
Format: Keep a Changelog. Versioning: Semantic Versioning.

## [Unreleased]

### Added

- Every tag gets a GitHub release whose notes are the version's CHANGELOG.md section, the engine level and the fleetlint version it needs, the digest pin and the commands that verify the signature, SBOM and provenance (`scripts/release-notes.sh`). The signed OCI artifact stays the only asset. Releases for v0.1.0 to v0.5.0 were created after the fact with the same notes.

### Changed

- **A publish step counts in any workflow not run from a pull request.** Every artifact kind is detected in tag-triggered workflows and in workflows run on push, release, schedule, by hand (`workflow_dispatch`), on `create`, `workflow_run`, `workflow_call` or `repository_dispatch`; `if:` guards on the step are not evaluated, since a step that is there publishes eventually. Before, ten kinds (release-files, npm, pypi, crate, gem, maven, nuget, pub, android-app, apple-app) only looked at tag-triggered workflows and missed release-please, semantic-release and changesets, which publish from a push to main. Rehearsals still do not count: `--dry-run` (and `-n` for cargo and pub), GoReleaser `build`, `check` and `--snapshot`, and Apple builds (`xcodebuild -exportArchive`, fastlane `gym`, `flutter build ipa`) outside a tag-triggered workflow.
- CI installs fleetlint v0.8.1 (0.8.0 is a source tag without binaries).

## [0.5.0] - 2026-10-21

Engine level 7: needs fleetlint 0.8.0 or newer (published artifacts as named instances, `locate:`; `combine:` and `artifact_workflows()` are gone).

### Changed

- **Published artifacts are instances with names.** `artifacts/image.yaml`, `helm-chart.yaml` and `oci-artifact.yaml` say with `instances:` what tells one artifact of the kind from another (an image reference, a chart path, a registry reference), so two images in one repository are checked one by one and the report names the one that fails. The release implementations read the current artifact's workflows as `artifact.workflows` instead of `artifact_workflows(kind)`.
- Rules whose finding is about one file say so with `locate:` (189 rules), a CEL string expression (`'"README.md"'`, or computed from the task runner or the publishing workflow), so reports and SARIF point at the file.
- **`combine:` is gone.** `changelog/generator-configured` is one rule at the root with every generator as a satisfier (the per-stack implementations are aliases of it). `release/reproducible-build` is keyed by artifact: `release-files` (Go `-trimpath` with `mod_timestamp`, Cargo `--locked` with `--remap-path-prefix`, .NET `ContinuousIntegrationBuild`, or `SOURCE_DATE_EPOCH`), `image` (`SOURCE_DATE_EPOCH` or `rewrite-timestamp=true`), `maven`, `nuget`, `pypi` and `gem`; `npm`, `crate` and `pub` are source packages and do not apply, `helm-chart`, `oci-artifact`, `android-app` and `apple-app` are not checked. The stack-keyed ids are aliases of the capability.

### Security

- The publish workflow passes the pushed digest to `cosign` through the environment instead of interpolating a step output into `run:`.

## [0.4.0] - 2026-10-21

Engine level 6: needs fleetlint 0.7.0 or newer (typed refs).

### Changed

- **Every rule names its source with a typed ref.** A `refs:` entry is a bare URL or `{url, kind, section}`; the kinds, strongest first, are `spec` (a language or format specification, or an official tool), `guidance` (official guidance such as Go Code Review Comments, PEP 8 or GitHub's own docs), `style-guide` (a public style guide), `standard` (a process or supply-chain standard: NIST SSDF, SLSA, OpenSSF Scorecard, Keep a Changelog, Conventional Commits, Semantic Versioning, EditorConfig, Reproducible Builds), `check` (a linter's check, `section` is its id) and `convention` (a convention of this catalog with no outside source, pointing at the fleetlint.org baseline page). Implementations carry the stack's own source (the tool's documentation), capabilities the general one; the fleetlint.org page links stay as untyped refs after the typed sources.

## [0.3.1] - 2026-10-21

### Added

- Artifact kind `oci-artifact`: a non-image artifact pushed to a registry with `oras push`, `flux push artifact`, `cosign upload blob` or `regctl artifact put`, in the workflow step or in a repository script the step runs. Implementations of the four release capabilities for it: signing by digest with cosign or Notation, an SBOM attached with `cosign attest` or `oras attach`, provenance with actions/attest-build-provenance, `cosign attest` or the SLSA generic generator; trusted publishing does not apply.
- The published catalog version carries an SPDX SBOM (`cosign attest`) and SLSA build provenance (actions/attest-build-provenance) on its digest, next to the signature.

### Changed

- `release/checksums` does not ask a release that publishes an OCI artifact or a chart for a checksums file: the registry digest is the checksum.
- The catalog repository checks itself at tier 1.
- The `image` kind sees GoReleaser `dockers_v2:` and `kos:` next to `dockers:`; `release/sbom/image` accepts `dockers_v2` with `sbom: true` and ko's SBOM, `release/provenance/image` accepts a `dockers_v2` build (buildx attaches provenance unless `provenance: false`).

## [0.3.0] - 2026-10-20

Engine level 5: needs fleetlint 0.6.0 or newer (artifacts as catalog data, capabilities keyed by artifact, `metadata.label`).

### Added

- **Artifacts are catalog data.** `artifacts/<id>.yaml` says what a release publishes and how it is detected: a CEL expression returning the workflows that publish it. Twelve kinds: `release-files`, `image`, `helm-chart`, `npm`, `pypi`, `crate`, `gem`, `maven`, `nuget`, `pub`, `android-app`, `apple-app`. A repository pins `facts: {artifacts: [...]}` when detection misses one.

### Changed

- **`release/signing`, `release/sbom`, `release/provenance` and `release/trusted-publishing` are checked per published artifact** (`by: artifact`), not per stack: a repository that publishes binaries and an npm package must sign both, and the report says which one is not. Each kind has its own implementation (`rules/release/<name>/<artifact>.yaml`); where a check means nothing for a kind (an SBOM for an npm package, a signature on a crate) the implementation says so and reports not applicable. `combine: any` is gone from these capabilities; the stack-keyed ids (`release/signing/node`, ...) are aliases of the new ones. Images pushed from a push- or schedule-triggered workflow count as published, not only from tag-triggered ones.
- `metadata.layer` is `metadata.label` (the old key is still read).

## [0.2.1] - 2026-10-15

The rules are those of 0.2.0.

### Fixed

- The publish workflow installs the fleetlint that reads engine level 4, so the published artifact is loaded the way a repository pins it before the run passes; 0.2.0 was pushed and signed, but that last check ran the old fleetlint.

## [0.2.0] - 2026-10-15

Engine level 4: needs fleetlint 0.5.0 or newer (`combine:` on capabilities, the stack-neutral implementation running where no stack is detected).

### Changed

- **More capabilities.** These rules state the requirement once and check it per stack at `rules/<family>/<name>/<stack>.yaml`; overrides, exceptions and policy on the old id reach every implementation:
  - `release/signing`: cosign, GoReleaser and gpg for every stack (`any`); the APK/AAB signature and PGP for Maven Central (kotlin), PGP and jarsigner (java), APK signing and Apple code signing (flutter), npm provenance (node), PyPI attestations and sigstore-python (python), `nuget sign` and Authenticode (dotnet), release-gem attestations and signed gems (ruby). A Kotlin or Android release no longer needs cosign.
  - `release/sbom`: syft, trivy, cdxgen, cyclonedx CLIs, sbom-tool, GoReleaser, BuildKit for every stack; the CycloneDX Gradle task and Maven plugin (kotlin, java; now only when the release runs them), `npm sbom` and cyclonedx-npm (node), cargo-cyclonedx (rust), CycloneDX .NET (dotnet), the CycloneDX Composer plugin (php).
  - `release/provenance`: attest action, SLSA generator, `cosign attest`/`attest-blob`, BuildKit, GitLab metadata for every stack; npm provenance (node), PyPI attestations (python), release-gem attestations (ruby).
  - `release/api-compatibility`: each stack decides when its project is a library and names its tool; new for java (japicmp, revapi), flutter (dart_apitool) and dotnet (package validation, ApiCompat), and kotlin also accepts the Kotlin Gradle plugin's `abiValidation`.
  - `quality/mutation-gate`: one tool per stack, new for dotnet (Stryker.NET), php (Infection), ruby (mutant) and cpp (Mull); not applicable to Hugo.
  - `quality/duplication-gate`: jscpd, Simian, duplo and Sonar for every stack; golangci `dupl` (go), pylint (python), PMD CPD (java, kotlin). Each stack's code needs a detector: a Go and Python repository that runs only pylint now fails for Go.
  - `public/manifest-license`: one manifest per stack, new for java and kotlin (POM licenses), ruby (gemspec), php (composer.json), dotnet (`PackageLicenseExpression`) and cpp (vcpkg.json, conanfile); not applicable to go, flutter and hugo, whose manifests have no license field.
  - `deps/forks-documented`: one manifest per stack, new for ruby (Gemfile `git:`/`github:`) and php (composer VCS repositories).
  - `release/trusted-publishing`: one registry per stack (npm, PyPI, crates.io, RubyGems), each with its own publish commands and token variables; not applicable to go, php, java, kotlin, hugo and cpp.
  - `ci/sast`: CodeQL, Semgrep, GitLab SAST and SonarQube for every stack; Brakeman (ruby) and Psalm taint analysis (php).
  - `changelog/generator-configured`: the generators' own files for every stack; changesets and semantic-release in package.json (node), git-cliff, commitizen and towncrier in pyproject.toml (python), git-cliff Cargo metadata (rust).
- `combine: any` on the release capabilities and `changelog/generator-configured`: a release and a changelog belong to the repository, so one stack's implementation passing is enough. The others keep `combine: all`.
- The stack-neutral implementations of `deps/vulnerability-scan`, `deps/license-check` and `quality/coverage-threshold` do not apply without a detected stack (no dependency manifest, no code); the release ones do, so a repository that only releases a container image is checked.
- `lint/configured/kotlin` requires detekt 2 (`dev.detekt`); a build that still applies detekt 1.x (`io.gitlab.arturbosch.detekt`, built on Kotlin 2.0, fails on JDK 25) gets a partial with that message. A version-catalog alias counts when `gradle/libs.versions.toml` maps it to `dev.detekt`. The fix names detekt 2 and a convention plugin for multi-module builds.
- The `kotlin/detekt.yml` template uses detekt 2's configuration keys (`allowedLines`, `allowedComplexity`, `allowedFunctionsPerFile`/`PerClass`, `UnusedPrivateFunction`, no `build.maxIssues`), which detekt 2 refused in the old template, and excludes every test source set (`commonTest`, `androidHostTest`, ...) from MagicNumber and FunctionNaming and `@Composable` functions from FunctionNaming, MagicNumber and LongParameterList.

### Fixed

- The check and release workflow templates install fleetlint from its release archive with a checksum check instead of `go install`, which needed Go on non-Go stacks.
- The Kotlin and Java workflow templates set `java-version: "25"` instead of reading a `.java-version` that `fleetlint fix` never writes.
- The pre-commit template leaves the generated `gradlew` script to Gradle (no shellcheck or shfmt rewrites) and lets `.bat`/`.cmd` files keep the CRLF that the `.gitattributes` template asks for.

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
