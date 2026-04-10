//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package python

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	cfgDep "github.com/ActiveMemory/ctx/internal/config/dep"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/io"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
)

// Name returns the ecosystem label.
//
// Returns:
//   - string: The Python ecosystem identifier
func (p *Builder) Name() string {
	return cfgDep.EcosystemPython
}

// Detect returns true if requirements.txt or pyproject.toml
// exists.
//
// Returns:
//   - bool: True when a Python manifest file is present
func (p *Builder) Detect() bool {
	for _, manifest := range []string{
		cfgDep.FileRequirements, cfgDep.FilePyproject,
	} {
		if _, err := os.Stat(manifest); err == nil {
			return true
		}
	}
	return false
}

// Build produces an adjacency list of Python dependencies.
//
// Parameters:
//   - external: If true, include dev dependencies
//
// Returns:
//   - map[string][]string: Dependency graph
//   - error: Non-nil if manifest parsing fails
func (p *Builder) Build(
	external bool,
) (map[string][]string, error) {
	if _, statErr := os.Stat(
		cfgDep.FilePyproject,
	); statErr == nil {
		return BuildPyprojectGraph(external)
	}
	return BuildRequirementsGraph()
}

// BuildRequirementsGraph parses requirements.txt and
// returns a flat dep list.
//
// Returns:
//   - map[string][]string: Dependency graph
//   - error: Non-nil if requirements.txt cannot be read
func BuildRequirementsGraph() (map[string][]string, error) {
	deps, parseErr := ParseRequirementsTxt(
		cfgDep.FileRequirements,
	)
	if parseErr != nil {
		return nil, parseErr
	}

	graph := make(map[string][]string)
	if len(deps) > 0 {
		sort.Strings(deps)
		graph[cfgDep.PyGraphRoot] = deps
	}
	return graph, nil
}

// ParseRequirementsTxt extracts package names from a
// requirements.txt file.
//
// Parameters:
//   - path: Path to requirements.txt
//
// Returns:
//   - []string: Package names
//   - error: Non-nil if file cannot be read
func ParseRequirementsTxt(
	path string,
) ([]string, error) {
	f, openErr := io.SafeOpenUserFile(path)
	if openErr != nil {
		return nil, openErr
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			ctxLog.Warn(warn.Close, path, closeErr)
		}
	}()

	var deps []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" ||
			strings.HasPrefix(line, token.PrefixComment) {
			continue
		}
		if strings.HasPrefix(
			line, cfgDep.TomlOptionPrefix,
		) {
			continue
		}
		name := ExtractPkgName(line)
		if name != "" {
			deps = append(deps, name)
		}
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return nil, scanErr
	}
	return deps, nil
}

// ExtractPkgName extracts the package name from a
// requirements line. Handles: package==1.0,
// package>=1.0, package[extra]>=1.0, package ; markers
//
// Parameters:
//   - line: Requirements line to parse
//
// Returns:
//   - string: Lowercase package name
func ExtractPkgName(line string) string {
	if idx := strings.Index(
		line, cfgDep.TomlComment,
	); idx >= 0 {
		line = line[:idx]
	}
	if idx := strings.Index(
		line, cfgDep.TomlSemicolon,
	); idx >= 0 {
		line = line[:idx]
	}
	line = strings.TrimSpace(line)

	if idx := strings.IndexAny(
		line, cfgDep.PyVersionSpecChars,
	); idx >= 0 {
		return strings.ToLower(
			strings.TrimSpace(line[:idx]),
		)
	}
	return strings.ToLower(strings.TrimSpace(line))
}

// BuildPyprojectGraph parses pyproject.toml for
// dependencies.
//
// Parameters:
//   - includeDevDeps: If true, include dev dependencies
//
// Returns:
//   - map[string][]string: Dependency graph
//   - error: Non-nil if pyproject.toml cannot be read
func BuildPyprojectGraph(
	includeDevDeps bool,
) (map[string][]string, error) {
	data, readErr := io.SafeReadUserFile(
		cfgDep.FilePyproject,
	)
	if readErr != nil {
		return nil, readErr
	}

	content := string(data)
	deps := ParsePyprojectDeps(content, cfgDep.PyDeps)

	if includeDevDeps {
		devDeps := ParsePyprojectDeps(
			content, cfgDep.PyDevDeps,
		)
		devDeps = append(
			devDeps,
			ParsePyprojectDeps(content, cfgDep.PyDev)...,
		)
		deps = append(deps, devDeps...)
	}

	seen := make(map[string]bool)
	var unique []string
	for _, d := range deps {
		if !seen[d] {
			seen[d] = true
			unique = append(unique, d)
		}
	}

	graph := make(map[string][]string)
	if len(unique) > 0 {
		sort.Strings(unique)
		graph[cfgDep.PyGraphRoot] = unique
	}
	return graph, nil
}

// ParsePyprojectDeps extracts dependency names from a TOML
// array section.
//
// Parameters:
//   - content: Full pyproject.toml content
//   - sectionSuffix: Section name suffix
//
// Returns:
//   - []string: Extracted dependency names
func ParsePyprojectDeps(
	content string, sectionSuffix string,
) []string {
	lines := strings.Split(content, token.NewlineLF)
	var deps []string
	inSection := false
	inArray := false

	targets := []string{
		fmt.Sprintf(
			cfgDep.TomlSectionFmt,
			cfgDep.TomlProjectPrefix, sectionSuffix,
		),
		fmt.Sprintf(
			cfgDep.TomlSectionFmt,
			cfgDep.TomlPoetryPrefix, sectionSuffix,
		),
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(
			trimmed, cfgDep.TomlSectionOpen,
		) {
			inSection = false
			inArray = false
			for _, target := range targets {
				if trimmed == target {
					inSection = true
					break
				}
			}
			continue
		}

		if !inSection {
			if !inArray {
				for _, key := range []string{
					sectionSuffix + cfgDep.TomlArrayAssign1,
					sectionSuffix + cfgDep.TomlArrayAssign2,
					sectionSuffix + cfgDep.TomlArrayAssign3,
					sectionSuffix + cfgDep.TomlArrayAssign4,
				} {
					if strings.Contains(trimmed, key) {
						inArray = true
						idx := strings.Index(
							trimmed, cfgDep.TomlSectionOpen,
						)
						rest := trimmed[idx+1:]
						deps = append(
							deps,
							ParseArrayItems(rest)...,
						)
						if strings.Contains(
							rest, token.CloseBracket,
						) {
							inArray = false
						}
						break
					}
				}
			} else {
				deps = append(
					deps, ParseArrayItems(trimmed)...,
				)
				if strings.Contains(
					trimmed, token.CloseBracket,
				) {
					inArray = false
				}
			}
			continue
		}

		if trimmed == "" ||
			strings.HasPrefix(trimmed, token.PrefixComment) {
			continue
		}

		if idx := strings.Index(
			trimmed, token.KeyValueSep,
		); idx > 0 {
			name := strings.TrimSpace(trimmed[:idx])
			lower := strings.ToLower(name)
			if lower == cfgDep.EcosystemPython ||
				cfgDep.PyMetaKeys[lower] {
				continue
			}
			deps = append(deps, strings.ToLower(name))
		}
	}
	return deps
}

// ParseArrayItems extracts package names from a TOML
// array line.
//
// Parameters:
//   - line: TOML array line to parse
//
// Returns:
//   - []string: Extracted package names
func ParseArrayItems(line string) []string {
	line = strings.ReplaceAll(
		line, token.CloseBracket, "",
	)
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	var deps []string
	for _, item := range strings.Split(
		line, token.Comma,
	) {
		item = strings.TrimSpace(item)
		item = strings.Trim(item, token.Quotes)
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		name := ExtractPkgName(item)
		if name != "" {
			deps = append(deps, name)
		}
	}
	return deps
}
