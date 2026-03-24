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
	"github.com/ActiveMemory/ctx/internal/cli/sync/core"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dep"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgSync "github.com/ActiveMemory/ctx/internal/config/sync"
	"github.com/ActiveMemory/ctx/internal/entity"
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
func CheckPackageFiles(ctx *entity.Context) []core.Action {
	var actions []core.Action

	for f, d := range dep.Packages {
		if _, err := os.Stat(f); err == nil {
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
				actions = append(actions, core.Action{
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
func CheckConfigFiles(ctx *entity.Context) []core.Action {
	var actions []core.Action

	for _, cfg := range lookup.ConfigPatterns() {
		matches, _ := filepath.Glob(cfg.Pattern)
		if len(matches) > 0 {
			// Check if CONVENTIONS.md mentions this
			var convContent string
			if f := ctx.File(cfgCtx.Convention); f != nil {
				convContent = strings.ToLower(string(f.Content))
			}

			keyword := strings.ToLower(strings.TrimPrefix(cfg.Pattern, "."))
			keyword = strings.TrimSuffix(keyword, "*")
			if convContent == "" || !strings.Contains(convContent, keyword) {
				actions = append(actions, core.Action{
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
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []Action: Suggested actions for undocumented directories
func CheckNewDirectories(ctx *entity.Context) []core.Action {
	var actions []core.Action

	// Get ARCHITECTURE.md content
	var archContent string
	if f := ctx.File(cfgCtx.Architecture); f != nil {
		archContent = strings.ToLower(string(f.Content))
	}

	// Scan top-level directories
	entries, err := os.ReadDir(".")
	if err != nil {
		return actions
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, ".") || cfgSync.SkipDirs[name] {
			continue
		}

		if cfgSync.ImportantDirs[name] && !strings.Contains(archContent, name) {
			actions = append(actions, core.Action{
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

	return actions
}
