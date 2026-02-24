//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.

package config

// Pattern represents a config file pattern and its documentation topic.
//
// Fields:
//   - Pattern: Glob pattern to match (e.g., ".eslintrc*")
//   - Topic: Documentation topic (e.g., "linting conventions")
type Pattern struct {
	Pattern string
	Topic   string
}

// Patterns lists config files that should be documented in CONVENTIONS.md.
//
// Used by sync to suggest documenting project configuration.
var Patterns = []Pattern{
	{".eslintrc*", "linting conventions"},
	{".prettierrc*", "formatting conventions"},
	{"tsconfig.json", "TypeScript configuration"},
	{".editorconfig", "editor configuration"},
	{"Makefile", "build commands"},
	{"Dockerfile", "containerization"},
}
