// Package catalog holds the fleetlint presets, rules, stacks and the
// templates their fixes write, as data. The fleetlint binary embeds one version of this module;
// repositories can also reference a preset file from this repository at a
// tag or commit.
package catalog

import "embed"

// Presets contains presets/<name>.yaml: minimal, recommended, strict, oss
// and the add-on slop. In fleetlint they are referenced as fleetlint:<name>.
// A preset selects rules from the library with `- use: <id or glob>`.
//
//go:embed presets/*.yaml
var Presets embed.FS

// Rules contains rules/<family>/<name>.yaml, one rule per file, the path
// spelling the rule's id. This is the library the presets select from.
//
//go:embed all:rules
var Rules embed.FS

// Stacks contains stacks/<id>.yaml: how a stack is detected and what its
// dev environment needs. Rules implement their capabilities per stack id.
//
//go:embed stacks/*.yaml
var Stacks embed.FS

// Templates contains templates/: the files `fleetlint fix` creates, plus the
// fragments it assembles the hook configuration and the check workflow from.
//
//go:embed all:templates
var Templates embed.FS
