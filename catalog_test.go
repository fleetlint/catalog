package catalog_test

import (
	"io/fs"
	"strings"
	"testing"

	"github.com/fleetlint/catalog"
)

// The engine that understands the rules lives in fleetlint, which tests every
// rule against fixtures when it takes a new version of this module. Here we
// only make sure the files it expects are present and not empty.
func TestPresetsAndTemplatesArePresent(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"minimal", "recommended", "slop"} {
		b, err := catalog.Presets.ReadFile("presets/" + name + ".yaml")
		if err != nil {
			t.Fatal(err)
		}
		text := string(b)
		if !strings.Contains(text, "apiVersion: fleetlint.org/v1") || !strings.Contains(text, "\n  - id: ") {
			t.Errorf("preset %s does not look like a catalog", name)
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
