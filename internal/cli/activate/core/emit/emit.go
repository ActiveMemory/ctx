//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package emit

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/env"
	cfgShell "github.com/ActiveMemory/ctx/internal/config/shell"
)

// emitters maps supported shell identifiers to their set-emitter.
// Unknown shells fall back to POSIX export semantics via the
// default branch in Set.
var emitters = map[string]emitter{
	cfgShell.Bash: posixSet,
	cfgShell.Zsh:  posixSet,
	cfgShell.Sh:   posixSet,
}

// unsetters maps supported shell identifiers to their unset-emitter.
// Unknown shells fall back to POSIX unset semantics via the default
// branch in Unset.
var unsetters = map[string]emitter{
	cfgShell.Bash: posixUnset,
	cfgShell.Zsh:  posixUnset,
	cfgShell.Sh:   posixUnset,
}

// DetectShell returns the shell identifier to emit for.
//
// Priority: explicit override > basename of $SHELL > bash fallback.
// The returned value is always lowercase and suitable as a key into
// the [emitters] / [unsetters] tables.
//
// Parameters:
//   - override: explicit --shell flag value ("" to auto-detect).
//
// Returns:
//   - string: one of [cfgShell.Bash], [cfgShell.Zsh], [cfgShell.Sh],
//     or the original override (callers treat unknowns as POSIX).
func DetectShell(override string) string {
	if override != "" {
		return strings.ToLower(override)
	}
	if s := os.Getenv(env.Shell); s != "" {
		return strings.ToLower(filepath.Base(s))
	}
	return cfgShell.Bash
}

// Set returns the shell command that exports CTX_DIR=path, ending
// with a newline so the output is directly consumable by
// `eval "$(ctx activate)"`.
//
// Parameters:
//   - shell: result of [DetectShell].
//   - path:  absolute path to the selected context directory.
//
// Returns:
//   - string: one-line export statement with trailing newline.
func Set(shell, path string) string {
	fn, ok := emitters[shell]
	if !ok {
		fn = posixSet
	}
	return fn(env.CtxDir, shellQuote(path))
}

// Unset returns the shell command that clears CTX_DIR for the
// current shell, ending with a newline.
//
// Parameters:
//   - shell: result of [DetectShell].
//
// Returns:
//   - string: one-line unset statement with trailing newline.
func Unset(shell string) string {
	fn, ok := unsetters[shell]
	if !ok {
		fn = posixUnset
	}
	return fn(env.CtxDir, "")
}
