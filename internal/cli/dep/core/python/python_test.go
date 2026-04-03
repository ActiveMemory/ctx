//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package python

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractPkgName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"requests", "requests"},
		{"requests==2.28.0", "requests"},
		{"requests>=2.28.0", "requests"},
		{"Flask>=2.0,<3.0", "flask"},
		{"uvicorn[standard]>=0.18.0", "uvicorn"},
		{"my-package~=1.0", "my-package"},
		{"Django ; python_version>='3.8'", "django"},
		{"boto3 # AWS SDK", "boto3"},
	}
	for _, tt := range tests {
		if got := ExtractPkgName(tt.input); got != tt.want {
			t.Errorf(
				"ExtractPkgName(%q) = %q, want %q",
				tt.input, got, tt.want,
			)
		}
	}
}

func TestBuilder_Requirements(t *testing.T) {
	orig, getErr := os.Getwd()
	if getErr != nil {
		t.Fatal(getErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	tmp := t.TempDir()
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}

	reqs := `# Core dependencies
flask>=2.0
requests==2.28.0
gunicorn

# Skip these
-r base.txt
--index-url https://pypi.org/simple/

# With extras
uvicorn[standard]>=0.18.0
`
	reqsPath := filepath.Join(tmp, "requirements.txt")
	if writeErr := os.WriteFile(
		reqsPath, []byte(reqs), 0o644,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	b := &Builder{}
	if !b.Detect() {
		t.Fatal(
			"Builder.Detect() = false with requirements.txt",
		)
	}

	graph, buildErr := b.Build(false)
	if buildErr != nil {
		t.Fatalf("Build(false) failed: %v", buildErr)
	}

	deps, ok := graph["project"]
	if !ok {
		t.Fatalf("expected 'project' key, got %v", graph)
	}
	if len(deps) != 4 {
		t.Errorf(
			"expected 4 deps, got %d: %v",
			len(deps), deps,
		)
	}
}

func TestBuilder_Pyproject(t *testing.T) {
	orig, getErr := os.Getwd()
	if getErr != nil {
		t.Fatal(getErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	tmp := t.TempDir()
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}

	toml := `[project]
name = "my-project"
version = "1.0.0"

[project.dependencies]
flask = ">=2.0"
requests = ">=2.28"

[project.dev-dependencies]
pytest = ">=7.0"
mypy = ">=1.0"
`
	tomlPath := filepath.Join(tmp, "pyproject.toml")
	if writeErr := os.WriteFile(
		tomlPath, []byte(toml), 0o644,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	b := &Builder{}
	if !b.Detect() {
		t.Fatal(
			"Builder.Detect() = false with pyproject.toml",
		)
	}

	graph, buildErr := b.Build(false)
	if buildErr != nil {
		t.Fatalf("Build(false) failed: %v", buildErr)
	}
	deps := graph["project"]
	if len(deps) != 2 {
		t.Errorf(
			"expected 2 deps without dev, got %d: %v",
			len(deps), deps,
		)
	}

	graphFull, buildErr := b.Build(true)
	if buildErr != nil {
		t.Fatalf("Build(true) failed: %v", buildErr)
	}
	depsFull := graphFull["project"]
	if len(depsFull) != 4 {
		t.Errorf(
			"expected 4 deps with dev, got %d: %v",
			len(depsFull), depsFull,
		)
	}
}

func TestBuilder_PyprojectInlineArray(t *testing.T) {
	orig, getErr := os.Getwd()
	if getErr != nil {
		t.Fatal(getErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	tmp := t.TempDir()
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}

	toml := `[project]
name = "my-project"
dependencies = [
    "requests>=2.28",
    "click>=8.0",
]
`
	tomlPath := filepath.Join(tmp, "pyproject.toml")
	if writeErr := os.WriteFile(
		tomlPath, []byte(toml), 0o644,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	b := &Builder{}
	graph, buildErr := b.Build(false)
	if buildErr != nil {
		t.Fatalf("Build(false) failed: %v", buildErr)
	}
	deps := graph["project"]
	if len(deps) != 2 {
		t.Errorf(
			"expected 2 deps, got %d: %v",
			len(deps), deps,
		)
	}
}

func TestBuilder_DetectPyproject(t *testing.T) {
	orig, getErr := os.Getwd()
	if getErr != nil {
		t.Fatal(getErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	tmp := t.TempDir()
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}

	minToml := "[project]\nname = \"test\"\n"
	minPath := filepath.Join(tmp, "pyproject.toml")
	if writeErr := os.WriteFile(
		minPath, []byte(minToml), 0o644,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	b := &Builder{}
	if !b.Detect() {
		t.Fatal(
			"Builder.Detect() = false with pyproject.toml",
		)
	}
}
