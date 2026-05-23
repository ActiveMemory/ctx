//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compliance

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// TestNoDirectStringsToLower bans `strings.ToLower(...)`
// calls in every .go file under the project tree except
// `internal/i18n/fold.go` itself. The replacement is
// `internal/i18n.Fold`, which is Unicode-correct for
// case-insensitive comparison.
//
// Why bytes-level lowercasing is wrong: strings.ToLower is
// locale-naive (Turkish İ→i̇ vs Turkish-aware i, German
// ß→ß vs ss, Greek final-sigma, etc.). Every use on
// potentially non-ASCII input is a latent i18n bug.
//
// Why no allowlist: the historical callsites (URL schemes,
// Go identifiers, ASCII keyword matching) happen to be
// safe because their input is ASCII-bounded, but
// `cases.Fold` produces byte-identical output for ASCII —
// so swapping to i18n.Fold is behavior-preserving on the
// safe paths and bug-fixing on the unsafe ones. A
// "grandfathered" allowlist would legitimize the broken
// windows; an annotation scheme would rot. The single
// structural exception is `internal/i18n/fold.go`, which
// is where the upstream `cases.Fold` call lives and
// cannot reasonably call its own helper (chicken-egg).
//
// See specs/i18n-fold-helper-and-ban.md.
func TestNoDirectStringsToLower(t *testing.T) {
	root := projectRoot(t)
	fset := token.NewFileSet()

	var violations []string
	for _, path := range allGoFiles(t, root) {
		// Structural exception: the i18n package's own
		// implementation is the only legitimate site of a
		// direct strings.ToLower equivalent — it calls the
		// upstream cases.Fold. The exception is matched on
		// the relative path under internal/i18n/ rather
		// than a hardcoded filename so renames within the
		// package don't silently re-enable bans.
		rel, _ := filepath.Rel(root, path)
		if strings.HasPrefix(filepath.ToSlash(rel), "internal/i18n/") {
			continue
		}

		file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			// Don't fail the whole test on a parse error
			// in one file — let go vet / go build surface
			// it. Skip and continue.
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if sel.Sel.Name != "ToLower" {
				return true
			}
			ident, ok := sel.X.(*ast.Ident)
			if !ok || ident.Name != "strings" {
				return true
			}
			pos := fset.Position(call.Pos())
			violations = append(violations,
				rel+":"+strconv.Itoa(pos.Line)+": direct strings.ToLower call — use internal/i18n.Fold instead",
			)
			return true
		})
	}

	for _, v := range violations {
		t.Error(v)
	}
}
