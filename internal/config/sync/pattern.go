//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

// Config file glob patterns checked by ctx sync.
const (
	PatternEslint     = ".eslintrc*"
	PatternPrettier   = ".prettierrc*"
	PatternTSConfig   = "tsconfig.json"
	PatternEditorConf = ".editorconfig"
	PatternMakefile   = "Makefile"
	PatternDockerfile = "Dockerfile"
)

// Action type constants for sync check results.
const (
	ActionDeps   = "DEPS"
	ActionConfig = "CONFIG"
	ActionNewDir = "NEW_DIR"
)

// Packages maps manifest filenames to their ecosystem descriptions.
// Used by sync to detect projects and suggest dependency documentation.
var Packages = map[string]string{
	"package.json":     "Node.js dependencies",
	"go.mod":           "Go module dependencies",
	"Cargo.toml":       "Rust dependencies",
	"requirements.txt": "Python dependencies",
	"Gemfile":          "Ruby dependencies",
}

// KeywordDependencies is the search keyword for dependency documentation.
const KeywordDependencies = "dependencies"

// ImportantDirs lists top-level directories that should be documented
// in ARCHITECTURE.md.
var ImportantDirs = map[string]bool{
	"api": true, "app": true, "cmd": true, "components": true,
	"internal": true, "lib": true, "pkg": true, "services": true,
	"src": true, "web": true,
}

// SkipDirs lists directories excluded from sync directory scanning.
var SkipDirs = map[string]bool{
	"build": true, "dist": true, "node_modules": true, "vendor": true,
}
