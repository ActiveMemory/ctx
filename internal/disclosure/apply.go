//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	"path/filepath"

	cfgFile "github.com/ActiveMemory/ctx/internal/config/file"
	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
)

// Apply is the mover: it digests a plan's staged entries into per-theme
// files and folds their gists into the root's ## Themes, under the spec's
// guards (specs/progressive-disclosure.md §Guards). It is the first pass
// that writes a canonical knowledge file.
//
// Order is load-bearing:
//  1. Refuse a non-knowledge file, an un-digestible kind (conventions are
//     a later milestone), or a malformed root (Validate) — before any
//     write.
//  2. Compute the lossless staging split and the theme-file appends.
//  3. Append entry bodies to their theme files (additive first).
//  4. Verify byte-presence by re-reading each theme file. Any miss aborts
//     with the root untouched.
//  5. Rewrite the root ONCE, last: remaining staging + folded gists. The
//     root write is the final syscall, so any earlier failure leaves the
//     root byte-identical (append→verify→remove; never remove-then-append).
//
// Fail-loud, no auto-repair. Worst-case crash duplicates a theme-file
// append (detectable, recoverable), never loses an entry.
//
// Parameters:
//   - rootPath: path to the canonical knowledge file (LEARNINGS/DECISIONS)
//   - plan: the digest plan (theme assignments + gists)
//   - ctxDir: the context directory theme files are written under
//
// Returns:
//   - ApplyResult: entries moved and theme slugs touched
//   - error: a disclosure sentinel, a path-bearing IO error, or nil
func Apply(rootPath string, plan Plan, ctxDir string) (ApplyResult, error) {
	kind, ok := KindFor(filepath.Base(rootPath))
	if !ok {
		return ApplyResult{}, errDisc.NotAKnowledgeFile(rootPath)
	}
	noun, ok := ThemeDir(kind)
	if !ok {
		return ApplyResult{}, errDisc.ErrApplyNotEntryKind
	}

	raw, readErr := internalIo.SafeReadUserFile(filepath.Clean(rootPath))
	if readErr != nil {
		return ApplyResult{}, readErr
	}
	root := Parse(string(raw), kind)
	if valErr := Validate(root); valErr != nil {
		return ApplyResult{}, valErr
	}

	moveIDs, planErr := FlattenPlan(plan)
	if planErr != nil {
		return ApplyResult{}, planErr
	}
	moved, remaining, splitErr := SplitStaging(root.Staging, moveIDs)
	if splitErr != nil {
		return ApplyResult{}, splitErr
	}

	// Empty plan: idempotent no-op, root untouched.
	if len(plan.Assignments) == 0 {
		return ApplyResult{}, nil
	}

	themeDir := filepath.Join(ctxDir, noun)
	if mkErr := internalIo.SafeMkdirAll(themeDir, cfgFs.PermExec); mkErr != nil {
		return ApplyResult{}, mkErr
	}

	// 3. Additive first: append each theme's bodies to its file.
	touched := make([]string, 0, len(plan.Assignments))
	for _, a := range plan.Assignments {
		path := filepath.Join(themeDir, a.Slug+cfgFile.ExtMarkdown)
		if appErr := appendTheme(path, a, moved); appErr != nil {
			return ApplyResult{}, appErr
		}
		touched = append(touched, a.Slug)
	}

	// 4. Verify byte-presence before touching the root.
	if vErr := verifyThemes(themeDir, plan, moved); vErr != nil {
		return ApplyResult{}, vErr
	}

	// 5. Rewrite the root once, last.
	newRoot := rewriteRoot(root, remaining, plan, noun)
	if wErr := internalIo.SafeWriteFile(
		rootPath, []byte(newRoot), cfgFs.PermFile,
	); wErr != nil {
		return ApplyResult{}, wErr
	}

	return ApplyResult{Moved: len(moveIDs), Themes: touched}, nil
}
