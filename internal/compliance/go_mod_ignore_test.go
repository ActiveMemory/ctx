//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compliance

import (
	"encoding/json"
	"os/exec"
	"testing"
)

// TestGoModIgnoresNodeModules verifies that the root go.mod
// carries an `ignore node_modules` directive.
//
// npm packages may ship Go sources (flatted does, under
// flatted/golang/). Without the directive, `./...` matches
// them as packages of this module, so `go build`, `go vet`,
// `go test`, and golangci-lint all compile and lint
// third-party code the moment anyone runs `npm install` in
// editors/vscode, ctx-desktop, or tools/typecheck/*. The
// bare form (no `./` prefix) matches a node_modules
// directory at any depth.
//
// The directive is read through `go mod edit -json` rather
// than by string search, so an ignore block or a reformat
// cannot fool the check.
//
// See specs/go-ignore-node-modules.md.
func TestGoModIgnoresNodeModules(t *testing.T) {
	root := projectRoot(t)

	cmd := exec.Command("go", "mod", "edit", "-json")
	cmd.Dir = root
	out, runErr := cmd.Output()
	if runErr != nil {
		t.Fatalf("go mod edit -json: %v", runErr)
	}

	var mod struct {
		Ignore []struct {
			Path string
		}
	}
	if parseErr := json.Unmarshal(out, &mod); parseErr != nil {
		t.Fatalf("parse go mod edit -json output: %v", parseErr)
	}

	for _, ig := range mod.Ignore {
		if ig.Path == "node_modules" {
			return
		}
	}
	t.Errorf("go.mod is missing `ignore node_modules`; " +
		"add it with: go mod edit -ignore=node_modules")
}
