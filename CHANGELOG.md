# Changelog

All notable changes to this project are documented here.
Format: Keep a Changelog. Versioning: Semantic Versioning.

## [Unreleased]

### Added

- The presets `minimal`, `recommended` and `slop` (94 rules) and the fix templates, moved here from the fleetlint repository. fleetlint embeds this module.
- Templates for a justfile and a Taskfile next to the Makefile, and workflow fragments that install those runners.

### Changed

- Rules that depend on a tool are outcome rules that list the known alternatives and accept more through `accept:`: `ci/check-workflow`, `lint/go-config`, `deps/vulnerability-scan`, `deps/license-check`, `repo/dev-environment`, and `taskrunner/container` (which replaces `taskrunner/devcontainer`).
- `taskrunner/targets` and `taskrunner/check-composition` read any task runner (make, just, task, package.json scripts) through `taskrunner.targets`, `taskrunner.deps` and `taskrunner.recipes`; they need a fleetlint that provides those.
