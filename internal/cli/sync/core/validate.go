//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dep"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
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
func CheckPackageFiles(ctx *entity.Context) []Action {
	var actions []Action

	for f, desc := range dep.Packages {
		if _, err := os.Stat(f); err == nil {
			// File exists, check if we have DEPENDENCIES.md or similar
			hasDepsDoc := false
			if f := ctx.File(ctxCfg.Dependency); f != nil {
				hasDepsDoc = true
			} else {
				for _, f := range ctx.Files {
					if strings.Contains(strings.ToLower(string(f.Content)),
						"dependencies",
					) {
						hasDepsDoc = true
						break
					}
				}
			}

			if !hasDepsDoc {
				actions = append(actions, Action{
					Type: "DEPS",
					File: ctxCfg.Architecture,
					Description: fmt.Sprintf(
						lookup.TextDesc(text.DescKeySyncDepsDescription),
						f, desc,
					),
					Suggestion: fmt.Sprintf(
						lookup.TextDesc(text.DescKeySyncDepsSuggestion),
						ctxCfg.Architecture, ctxCfg.Dependency,
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
			if f := ctx.File(ctxCfg.Convention); f != nil {
				convContent = strings.ToLower(string(f.Content))
			}

			keyword := strings.ToLower(strings.TrimPrefix(cfg.Pattern, "."))
			keyword = strings.TrimSuffix(keyword, "*")
			if convContent == "" || !strings.Contains(convContent, keyword) {
				actions = append(actions, Action{
					Type: "CONFIG",
					File: ctxCfg.Convention,
					Description: fmt.Sprintf(
						desc.TextDesc(text.DescKeySyncConfigDescription),
						matches[0], cfg.Topic,
					),
					Suggestion: fmt.Sprintf(
						desc.TextDesc(text.DescKeySyncConfigSuggestion),
						cfg.Topic, ctxCfg.Convention,
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
func CheckNewDirectories(ctx *entity.Context) []Action {
	var actions []Action

	// Get ARCHITECTURE.md content
	var archContent string
	if f := ctx.File(ctxCfg.Architecture); f != nil {
		archContent = strings.ToLower(string(f.Content))
	}

	// Scan top-level directories
	entries, err := os.ReadDir(".")
	if err != nil {
		return actions
	}

	importantDirs := []string{
		"src", "lib", "pkg", "internal",
		"cmd", "api", "web", "app", "services", "components",
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Skip hidden directories and common non-code directories
		if strings.HasPrefix(name, ".") ||
			name == "node_modules" ||
			name == "vendor" ||
			name == "dist" ||
			name == "build" {
			continue
		}

		// Check if this is an important directory not mentioned in ARCHITECTURE.md
		isImportant := false
		for _, imp := range importantDirs {
			if name == imp {
				isImportant = true
				break
			}
		}

		if isImportant && !strings.Contains(archContent, name) {
			actions = append(actions, Action{
				Type: "NEW_DIR",
				File: ctxCfg.Architecture,
				Description: fmt.Sprintf(
					desc.TextDesc(text.DescKeySyncDirDescription),
					name,
				),
				Suggestion: fmt.Sprintf(
					desc.TextDesc(text.DescKeySyncDirSuggestion),
					name, ctxCfg.Architecture,
				),
			})
		}
	}

	return actions
}
