package tddcheck_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lwmacct/260622-go-pkg-tddcheck/pkg/tddcheck"
)

func TestRules(t *testing.T) {
	tddcheck.Project{
		Root:   ".",
		Config: projectRules(),
	}.Assert(t)
}

func TestWriteProjectDoc(t *testing.T) {
	tddcheck.Project{
		Root:   ".",
		Config: projectRules(),
	}.WriteDoc(t, "")
}

func TestPublicAuthPackageBoundary(t *testing.T) {
	root := projectRoot(t)
	for _, file := range goFiles(t, filepath.Join(root, "pkg", "auth")) {
		imports := importsInFile(t, file)
		for _, importPath := range imports {
			if strings.Contains(importPath, "/internal/appcmd") ||
				strings.Contains(importPath, "/internal/config") ||
				strings.Contains(importPath, "/internal/infra/database") ||
				strings.Contains(importPath, "github.com/urfave/cli") ||
				strings.Contains(importPath, "251207-go-pkg-cfgm") ||
				strings.Contains(importPath, "260614-go-pkg-tlsreload") {
				t.Fatalf("pkg/auth must not depend on example app infrastructure: %s imports %s", file, importPath)
			}
		}
	}
}

func projectRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("project root not found")
		}
		dir = parent
	}
}

func projectRules() tddcheck.Config {
	cfg := tddcheck.DefaultConfig()
	cfg.DependencyLayerDirs = []string{"appcmd", "config", "handler", "infra", "repository", "service"}
	cfg.LayerRules = append(cfg.LayerRules,
		tddcheck.LayerDependencyRule{SourceLayer: "appcmd", TargetLayer: "handler", Message: "appcmd must compose through pkg/auth, not internal handler"},
		tddcheck.LayerDependencyRule{SourceLayer: "appcmd", TargetLayer: "service", Message: "appcmd must compose through pkg/auth, not internal service"},
		tddcheck.LayerDependencyRule{SourceLayer: "appcmd", TargetLayer: "repository", Message: "appcmd must compose through pkg/auth, not internal repository"},
		tddcheck.LayerDependencyRule{SourceLayer: "config", TargetLayer: "handler", Message: "config must not import handler"},
		tddcheck.LayerDependencyRule{SourceLayer: "config", TargetLayer: "service", Message: "config must not import service"},
		tddcheck.LayerDependencyRule{SourceLayer: "config", TargetLayer: "repository", Message: "config must not import repository"},
	)
	return cfg
}

func goFiles(t *testing.T, root string) []string {
	t.Helper()

	var files []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func importsInFile(t *testing.T, file string) []string {
	t.Helper()

	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, file, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatal(err)
	}
	imports := make([]string, 0, len(parsed.Imports))
	for _, item := range parsed.Imports {
		imports = append(imports, strings.Trim(item.Path.Value, `"`))
	}
	return imports
}
