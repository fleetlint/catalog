# fleetlint catalog

The rules [fleetlint](https://github.com/fleetlint/fleetlint) checks a repository against, and the files its fixes write, as data.

- `presets/minimal.yaml`: the floor for any repository.
- `presets/recommended.yaml`: the default; includes `minimal`.
- `presets/slop.yaml`: an add-on with heuristics for the traces careless or generated code leaves behind.
- `templates/`: what `fleetlint fix` creates (editor and git attributes, linter configuration, changelog, security policy, release workflows per stack, a Makefile, justfile and Taskfile with the task-runner contract) and the fragments it assembles the hook configuration and the check workflow from.

What the rules stand for is described in [the baseline](https://fleetlint.org/baseline/), [the stack profiles](https://fleetlint.org/stacks/) and [the slop catalog](https://fleetlint.org/slop/); the rule format in [writing rules](https://fleetlint.org/writing-rules/). The source of those pages is [fleetlint/docs](https://github.com/fleetlint/docs).

## Quickstart

Every fleetlint binary embeds one version of this catalog, so nothing is needed for the usual case:

```yaml
# .fleetlint.yaml
version: 1
extends: [fleetlint:recommended]
```

To use another version of the catalog than your binary ships, pin it; rules and templates then come from that version, fetched once and cached:

```yaml
version: 1
catalog: { version: v0.2.0 }     # a tag or a full commit SHA of this repository
extends: [fleetlint:recommended]
```

Each preset states the engine level it was written for (`metadata.engine`). A fleetlint that cannot run a version says so in one message instead of failing rule by rule.

To build your own baseline on these presets, write a catalog that includes one and adjusts it with `overrides:` (see [organisations](https://fleetlint.org/organisations/)).

## Changing a rule

Rules are YAML with [CEL](https://cel.dev) predicates. This repository only checks that the files are present; the engine that evaluates them is in fleetlint, which tests every rule against a passing and a failing fixture. A rule change therefore has two parts: the change here, and in fleetlint the fixture and the new catalog version in `go.mod`; and in fleetlint/docs the regenerated rule pages.

```sh
make check    # format, vet, tests
```

## Layout

`catalog.go` exposes `presets/` and `templates/` as embedded file systems (`catalog.Presets`, `catalog.Templates`); fleetlint imports this module for them.

## License

MIT, see `LICENSE`.

Files that fleetlint writes into your repository from these templates belong to that repository. Use, change and distribute them under whatever terms you like; no attribution or license notice is required for them.
