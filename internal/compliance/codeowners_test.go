//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compliance

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCodeownersSingleRule verifies that CODEOWNERS holds exactly
// one rule, `*`, naming every maintainer.
//
// GitHub applies only the last CODEOWNERS pattern that matches a
// file. Appending `* @new-maintainer` on its own line therefore
// replaces the previous owners instead of adding one, and any
// narrower pattern (`/internal/ @someone`) strips everyone else's
// ownership of that path. Both happened: five separate `*` lines
// left only the last name as owner of most of the repo, and eight
// path rules left a single owner of cmd/, internal/, docs/, specs/,
// hack/, .github/, and the Makefile.
//
// The policy is full ownership for every maintainer, so the file
// must be one `*` line. Owner names are not pinned here; adding a
// maintainer means appending a handle to that line.
//
// See specs/codeowners-full-ownership.md.
func TestCodeownersSingleRule(t *testing.T) {
	root := projectRoot(t)

	//nolint:gosec // constructed from test constants
	f, openErr := os.Open(filepath.Clean(filepath.Join(root, "CODEOWNERS")))
	if openErr != nil {
		t.Fatalf("open CODEOWNERS: %v", openErr)
	}
	defer func() { _ = f.Close() }()

	var rules []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		rules = append(rules, line)
	}
	if scanErr := scanner.Err(); scanErr != nil {
		t.Fatalf("read CODEOWNERS: %v", scanErr)
	}

	if len(rules) != 1 {
		t.Fatalf("CODEOWNERS has %d rules, want exactly 1 (`* @a @b ...`):\n\t%s",
			len(rules), strings.Join(rules, "\n\t"))
	}

	fields := strings.Fields(rules[0])
	if fields[0] != "*" {
		t.Errorf("CODEOWNERS rule pattern is %q, want %q", fields[0], "*")
	}
	if len(fields) < 2 {
		t.Errorf("CODEOWNERS `*` rule names no owners")
	}
	for _, owner := range fields[1:] {
		if !strings.HasPrefix(owner, "@") {
			t.Errorf("CODEOWNERS owner %q is not an @handle", owner)
		}
	}
}
