package catalog_test

import (
	"io/fs"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"

	"github.com/fleetlint/catalog"
)

// The engine that understands the rules lives in fleetlint, which tests every
// rule against fixtures when it takes a new version of this module. Here we
// only make sure the files it expects are present and not empty.
func TestPresetsAndTemplatesArePresent(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"minimal", "recommended", "slop", "strict", "oss"} {
		b, err := catalog.Presets.ReadFile("presets/" + name + ".yaml")
		if err != nil {
			t.Fatal(err)
		}
		text := string(b)
		selects := strings.Contains(text, "\n  - use: ") || strings.Contains(text, "\noverrides:\n")
		if !strings.Contains(text, "apiVersion: fleetlint.org/v1") || !selects {
			t.Errorf("preset %s does not look like a catalog that selects or overrides rules", name)
		}
	}
	files := 0
	err := fs.WalkDir(catalog.Templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		files++
		if b, rerr := catalog.Templates.ReadFile(path); rerr != nil || len(b) == 0 {
			t.Errorf("template %s is empty or unreadable", path)
		}
		return nil
	})
	if err != nil || files == 0 {
		t.Fatalf("templates: %d files, err=%v", files, err)
	}
	for _, stack := range []string{"go", "python", "rust", "node", "kotlin", "flutter"} {
		for _, name := range []string{"pre-commit/" + stack + ".yaml", "check/" + stack + ".yml"} {
			if _, err := catalog.Templates.ReadFile("templates/" + name); err != nil {
				t.Errorf("stack %s lacks %s", stack, name)
			}
		}
	}
	for _, name := range []string{"head.yml", "verify.yml", "goreleaser.yml", "dist.yml"} {
		if _, err := catalog.Templates.ReadFile("templates/release/" + name); err != nil {
			t.Errorf("release fragment %s missing", name)
		}
	}
}

// Every library file is one rule whose id spells its path, and every rule
// is selected by some preset: a rule nobody uses is dead data.
func TestLibraryRulesAreSelected(t *testing.T) {
	t.Parallel()
	selected := map[string]bool{}
	presets, err := fs.Glob(catalog.Presets, "presets/*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var globs []string
	for _, p := range presets {
		b, _ := catalog.Presets.ReadFile(p)
		var doc struct {
			Rules []struct {
				Use string `yaml:"use"`
			} `yaml:"rules"`
		}
		if err := yaml.Unmarshal(b, &doc); err != nil {
			t.Fatalf("%s: %v", p, err)
		}
		for _, r := range doc.Rules {
			if strings.Contains(r.Use, "*") {
				globs = append(globs, strings.TrimSuffix(r.Use, "*"))
			} else if r.Use != "" {
				selected[r.Use] = true
			}
		}
	}
	err = fs.WalkDir(catalog.Rules, "rules", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		id := strings.TrimSuffix(strings.TrimPrefix(p, "rules/"), ".yaml")
		b, _ := catalog.Rules.ReadFile(p)
		var rule struct {
			ID string `yaml:"id"`
		}
		if err := yaml.Unmarshal(b, &rule); err != nil || rule.ID != id {
			t.Errorf("%s: id %q must spell the path (%v)", p, rule.ID, err)
		}
		if selected[id] {
			return nil
		}
		for _, g := range globs {
			if strings.HasPrefix(id, g) {
				return nil
			}
		}
		t.Errorf("%s: no preset selects this rule", id)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
