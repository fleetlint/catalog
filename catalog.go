// Package catalog holds the fleetlint presets and the templates their fixes
// write, as data. The fleetlint binary embeds one version of this module;
// repositories can also reference a preset file from this repository at a
// tag or commit.
package catalog

import "embed"

// Presets contains presets/<name>.yaml: minimal, recommended and the add-on
// slop. In fleetlint they are referenced as fleetlint:<name>.
//
//go:embed presets/*.yaml
var Presets embed.FS

// Templates contains templates/: the files `fleetlint fix` creates, plus the
// fragments it assembles the hook configuration and the check workflow from.
//
//go:embed all:templates
var Templates embed.FS
