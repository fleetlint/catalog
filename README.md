# fleetlint catalog

The rules [fleetlint](https://github.com/fleetlint/fleetlint) checks a repository against, and the files its fixes write, as data.

- `presets/minimal.yaml`: the floor for any repository.
- `presets/recommended.yaml`: the default; includes `minimal`.
- `presets/slop.yaml`: an add-on with heuristics for the traces careless or generated code leaves behind.
- `templates/`: what `fleetlint fix` creates (editor and git attributes, linter configuration, changelog, security policy, release workflows per stack, a Makefile, justfile and Taskfile with the task-runner contract) and the fragments it assembles the hook configuration and the check workflow from.

What the rules stand for is described in fleetlint's `docs/baseline.md`, `docs/stacks.md` and `docs/slop.md`; the rule format in `docs/writing-rules.md`.

## Quickstart

Every fleetlint binary embeds one version of this catalog, so nothing is needed for the usual case:

```yaml
# .fleetlint.yaml
version: 1
extends: [fleetlint:recommended]
```

To use a different version of a preset than your binary ships, reference the file in this repository at a tag or commit:

```yaml
extends:
  - git+https://github.com/fleetlint/catalog.git//presets/recommended.yaml@<commit>
```

A preset loaded this way is a remote catalog: its rules apply, and fixes take their templates from the binary.

## Changing a rule

Rules are YAML with [CEL](https://cel.dev) predicates. This repository only checks that the files are present; the engine that evaluates them is in fleetlint, which tests every rule against a passing and a failing fixture. A rule change therefore has two parts: the change here, and in fleetlint the fixture, the new catalog version in `go.mod` and the regenerated rule pages (`make docs`).

```sh
make check    # format, vet, tests
```

## Layout

`catalog.go` exposes `presets/` and `templates/` as embedded file systems (`catalog.Presets`, `catalog.Templates`); fleetlint imports this module for them.

## License

MIT, see `LICENSE`.

Files that fleetlint writes into your repository from these templates belong to that repository. Use, change and distribute them under whatever terms you like; no attribution or license notice is required for them.
