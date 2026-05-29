//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compliance

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// ctxctlImportPrefix is the package-path prefix for the
// maintainer-only audit logic. The shipped ctx binary must
// never reach it.
const ctxctlImportPrefix = "github.com/ActiveMemory/ctx/internal/ctxctl"

// TestCtxBinaryExcludesCtxctl asserts that the shipped ctx
// binary's transitive import graph contains no
// internal/ctxctl package. This is the second of the two
// isolation layers from specs/ctxctl-bootstrap.md: the first
// is the module graph (ctx's go.mod does not require
// tools/ctxctl); this test guards against a regression that
// would re-introduce an import edge from cmd/ctx into the
// audit subtree.
func TestCtxBinaryExcludesCtxctl(t *testing.T) {
	root := projectRoot(t)

	cmd := exec.Command("go", "list", "-deps", "./cmd/ctx")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go list -deps ./cmd/ctx: %v", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		pkg := strings.TrimSpace(line)
		if pkg == "" {
			continue
		}
		if pkg == ctxctlImportPrefix ||
			strings.HasPrefix(pkg, ctxctlImportPrefix+"/") {
			t.Errorf(
				"cmd/ctx transitively imports %q; the audit "+
					"subtree must stay out of the shipped binary",
				pkg,
			)
		}
	}
}

// TestShippedHooksExcludeCheckAudit asserts the shipped
// hooks.json (installed by `ctx setup`) wires no check-audit
// hook. The audit channel is maintainer-only; taxing every
// end user's every prompt with an audit relay they have no
// producer for is exactly what the ctxctl migration removed.
func TestShippedHooksExcludeCheckAudit(t *testing.T) {
	root := projectRoot(t)
	hooksPath := filepath.Join(
		root, "internal", "assets", "claude", "hooks", "hooks.json",
	)

	data, err := os.ReadFile(filepath.Clean(hooksPath))
	if err != nil {
		t.Fatalf("read shipped hooks.json: %v", err)
	}

	if strings.Contains(string(data), "check-audit") {
		t.Errorf(
			"shipped hooks.json contains \"check-audit\"; the " +
				"audit relay must not ship to end users",
		)
	}
}
