// Package catalog holds the fleetlint presets, rules, stacks, artifacts and
// the templates their fixes write, as data. The fleetlint binary embeds one
// version of this module; a repository runs another with
// `catalog: {version: <tag or commit>}`.
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

// Artifacts contains artifacts/<id>.yaml: what a release publishes (files
// on a forge release, a container image, an npm package, ...) and how it is
// detected from the workflows that publish it. Release capabilities implement
// their checks per artifact id.
//
//go:embed artifacts/*.yaml
var Artifacts embed.FS

// Aliases is aliases.yaml: renamed rule ids and the ids that replace them.
//
//go:embed aliases.yaml
var Aliases []byte

// Templates contains templates/: the files `fleetlint fix` creates, plus the
// fragments it assembles the hook configuration and the check workflow from.
//
//go:embed all:templates
var Templates embed.FS
