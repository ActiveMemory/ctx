//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgSync "github.com/ActiveMemory/ctx/internal/config/sync"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// CheckPackageFiles detects package manager files without dependency
// documentation.
//
// Checks for common package files (package.json, go.mod, etc.) and suggests
// documenting dependencies if no DEPENDENCIES.md exists or ARCHITECTURE.md
// doesn't mention dependencies.
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []Action: Suggested actions for undocumented dependencies
func CheckPackageFiles(ctx *entity.Context) []Action {
	var actions []Action

	for f, d := range cfgSync.Packages {
		if _, statErr := os.Stat(f); statErr == nil {
			// File exists, check if we have DEPENDENCIES.md or similar
			hasDepsDoc := false
			if f := ctx.File(cfgCtx.Dependency); f != nil {
				hasDepsDoc = true
			} else {
				for _, f := range ctx.Files {
					if strings.Contains(strings.ToLower(string(f.Content)),
						cfgSync.KeywordDependencies,
					) {
						hasDepsDoc = true
						break
					}
				}
			}

			if !hasDepsDoc {
				actions = append(actions, Action{
					Type: cfgSync.ActionDeps,
					File: cfgCtx.Architecture,
					Description: fmt.Sprintf(
						lookup.TextDesc(text.DescKeySyncDepsDescription),
						f, d,
					),
					Suggestion: fmt.Sprintf(
						lookup.TextDesc(text.DescKeySyncDepsSuggestion),
						cfgCtx.Architecture, cfgCtx.Dependency,
					),
				})
			}
		}
	}

	return actions
}

// CheckConfigFiles detects config files not documented in CONVENTIONS.md.
//
// Scans for common configuration files (.eslintrc, .prettierrc, tsconfig.json,
// etc.) and suggests documenting them if CONVENTIONS.md doesn't mention the
// related topic.
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []Action: Suggested actions for undocumented configurations
func CheckConfigFiles(ctx *entity.Context) []Action {
	var actions []Action

	for _, cfg := range lookup.ConfigPatterns() {
		matches, _ := filepath.Glob(cfg.Pattern)
		if len(matches) > 0 {
			// Check if CONVENTIONS.md mentions this
			var convContent string
			if f := ctx.File(cfgCtx.Convention); f != nil {
				convContent = strings.ToLower(string(f.Content))
			}

			keyword := strings.ToLower(strings.TrimPrefix(cfg.Pattern, token.PrefixDot))
			keyword = strings.TrimSuffix(keyword, token.GlobStar)
			if convContent == "" || !strings.Contains(convContent, keyword) {
				actions = append(actions, Action{
					Type: cfgSync.ActionConfig,
					File: cfgCtx.Convention,
					Description: fmt.Sprintf(
						desc.Text(text.DescKeySyncConfigDescription),
						matches[0], cfg.Topic,
					),
					Suggestion: fmt.Sprintf(
						desc.Text(text.DescKeySyncConfigSuggestion),
						cfg.Topic, cfgCtx.Convention,
					),
				})
			}
		}
	}

	return actions
}

// CheckNewDirectories detects important directories not in ARCHITECTURE.md.
//
// Scans top-level directories for common code directories (src, lib, cmd, etc.)
// and suggests documenting them if ARCHITECTURE.md doesn't mention them.
// Skips hidden directories and common non-code directories (node_modules,
// vendor, dist, build).
//
// Returns (nil, nil) when the context directory is not declared: there is
// no project root to scan, which is an ordinary "nothing to suggest" state
// rather than an error. A resolver failure for any other reason, and a
// directory read failure, are propagated so the caller does not report a
// confident empty suggestion list when we actually failed to look.
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []Action: Suggested actions for undocumented directories
//   - error: non-nil on resolver failure (other than not-declared) or
//     on a directory read failure at the project root
func CheckNewDirectories(ctx *entity.Context) ([]Action, error) {
	var actions []Action

	// Get ARCHITECTURE.md content
	var archContent string
	if f := ctx.File(cfgCtx.Architecture); f != nil {
		archContent = strings.ToLower(string(f.Content))
	}

	// Scan top-level directories at the project root (parent of the
	// declared context directory). Under the explicit-context-dir
	// model this is authoritative; CWD may be a subdir.
	ctxDir, ctxErr := rc.ContextDir()
	if ctxErr != nil {
		return nil, ctxErr
	}
	projectRoot := filepath.Dir(ctxDir)
	entries, readDirErr := os.ReadDir(projectRoot)
	if readDirErr != nil {
		return nil, readDirErr
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, token.PrefixDot) || cfgSync.SkipDirs[name] {
			continue
		}

		if cfgSync.ImportantDirs[name] && !strings.Contains(archContent, name) {
			actions = append(actions, Action{
				Type: cfgSync.ActionNewDir,
				File: cfgCtx.Architecture,
				Description: fmt.Sprintf(
					desc.Text(text.DescKeySyncDirDescription),
					name,
				),
				Suggestion: fmt.Sprintf(
					desc.Text(text.DescKeySyncDirSuggestion),
					name, cfgCtx.Architecture,
				),
			})
		}
	}

	return actions, nil
}
