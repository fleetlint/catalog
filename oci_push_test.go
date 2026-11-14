package catalog_test

import (
	"go/build"
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// scripts/oci-push.sh publishes what the module embeds, no less: a pinned
// `catalog: {oci: ...}` reads only what the artifact holds.
func TestOCIPushPublishesWhatTheModuleEmbeds(t *testing.T) {
	t.Parallel()
	pkg, err := build.ImportDir(".", 0)
	if err != nil {
		t.Fatal(err)
	}
	var embedded []string
	for _, p := range pkg.EmbedPatterns {
		top, _, _ := strings.Cut(strings.TrimPrefix(p, "all:"), "/")
		embedded = append(embedded, top)
	}
	script, err := os.ReadFile("scripts/oci-push.sh")
	if err != nil {
		t.Fatal(err)
	}
	files := regexp.MustCompile(`(?m)^files=\(([^)]*)\)`).FindSubmatch(script)
	dirs := regexp.MustCompile(`find ([\w ]+?) -type f`).FindSubmatch(script)
	if files == nil || dirs == nil {
		t.Fatal("scripts/oci-push.sh no longer lists its files as files=(...) and find <dirs> -type f")
	}
	pushed := append(strings.Fields(string(files[1])), strings.Fields(string(dirs[1]))...)
	slices.Sort(embedded)
	slices.Sort(pushed)
	if !slices.Equal(slices.Compact(embedded), slices.Compact(pushed)) {
		t.Errorf("the module embeds %v, scripts/oci-push.sh pushes %v", embedded, pushed)
	}
}
