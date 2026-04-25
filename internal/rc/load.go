//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/config/env"
	"github.com/ActiveMemory/ctx/internal/config/file"
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	writeRC "github.com/ActiveMemory/ctx/internal/write/rc"
)

// load builds the runtime configuration under the
// single-source-anchor model
// (spec: specs/single-source-context-anchor.md).
//
// Lookup rules:
//
//   - When a context directory has been declared via CTX_DIR,
//     `.ctxrc` is read from
//     `filepath.Dir(ContextDir()) + "/.ctxrc"`: the project root,
//     which by contract is the parent of [ContextDir]. CWD has no
//     say. This is the "configuration belongs to the project root"
//     rule.
//   - When no context directory is declared, `.ctxrc` is not read
//     at all: there is no project to configure. Defaults apply.
//   - Environment overrides (CTX_TOKEN_BUDGET) are applied after the
//     YAML merge so users can tune per-session without editing the
//     file.
//
// Returns:
//   - *CtxRC: Configuration with file values (when .ctxrc is
//     readable) and environment overrides applied.
func load() *CtxRC {
	cfg := Default()

	rcPath, pathErr := ctxrcPath()
	switch {
	case pathErr == nil:
		data, readErr := ctxIo.SafeReadUserFile(rcPath)
		if readErr == nil {
			if yamlErr := yaml.Unmarshal(data, cfg); yamlErr != nil {
				writeRC.ParseWarning(rcPath, yamlErr)
			}
		}
	case errors.Is(pathErr, errCtx.ErrDirNotDeclared):
		// CTX_DIR not declared. **Expected** for exempt commands
		// (ctx init, activate, deactivate, doctor, version,
		// hub *, etc.) that legitimately call accessors before
		// any project exists; defaults are the right answer for
		// them. **Unexpected** for operating commands, which
		// should have been gated by [bootstrap/cmd.go]'s
		// PersistentPreRunE call to RequireContextDir before
		// reaching any RC accessor.
		//
		// If an operating command ever slips past that gate, this
		// branch would silently hand back default config
		// (token_budget = 8000, auto_archive = true, etc.) and
		// the user's .ctxrc settings would be invisibly ignored.
		// Emit a stderr breadcrumb so the silence is visible:
		// loud enough to surface during a missed-gate regression
		// in dev / CI, quiet enough to ignore in legitimate
		// exempt flows. Defaults still apply so the command can
		// keep running.
		logWarn.Warn(cfgWarn.RCNoContextDir)
	default:
		// Unexpected resolver failure (relative path,
		// non-canonical basename, etc.). Surface loudly rather
		// than swallowing; defaults still apply so commands that
		// do not require a project can still boot. Same noisy-TUI
		// principle documented on resolve.DirLine /
		// resolve.AppendDir.
		logWarn.Warn(cfgWarn.ContextDirResolve, pathErr)
	}

	if envBudget := os.Getenv(env.CtxTokenBudget); envBudget != "" {
		budget, parseErr := strconv.Atoi(envBudget)
		if parseErr == nil && budget > 0 {
			cfg.TokenBudget = budget
		}
	}

	return cfg
}

// ctxrcPath returns the absolute path to the `.ctxrc` file adjacent
// to the declared context directory.
//
// Returns:
//   - string: Absolute path to .ctxrc on success; "" on error.
//   - error: errCtx.ErrDirNotDeclared when no context directory has
//     been declared; any other resolver error from ContextDir is
//     propagated unchanged so the caller decides policy rather than
//     this helper silently returning an empty path.
func ctxrcPath() (string, error) {
	ctxDir, err := ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(ctxDir), file.CtxRC), nil
}
