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

// Patterns returns all config file glob patterns in detection order.
func Patterns() []string {
	return []string{
		PatternEslint,
		PatternPrettier,
		PatternTSConfig,
		PatternEditorConf,
		PatternMakefile,
		PatternDockerfile,
	}
}
