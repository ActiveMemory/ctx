//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deps

import (
	"bufio"
	"github.com/ActiveMemory/ctx/internal/config"
	"os"
	"sort"
	"strings"
)

const pythonEcosystem = "python"

// pythonBuilder implements GraphBuilder for Python projects.
type pythonBuilder struct{}

func (p *pythonBuilder) Name() string { return pythonEcosystem }

func (p *pythonBuilder) Detect() bool {
	for _, manifest := range []string{"requirements.txt", "pyproject.toml"} {
		if _, err := os.Stat(manifest); err == nil {
			return true
		}
	}
	return false
}

func (p *pythonBuilder) Build(external bool) (map[string][]string, error) {
	// Python builder always shows external deps — there's no internal
	// package graph without import tracing. The external flag controls
	// whether we include dev dependencies from pyproject.toml.
	if _, statErr := os.Stat("pyproject.toml"); statErr == nil {
		return buildPyprojectGraph(external)
	}
	return buildRequirementsGraph()
}

// buildRequirementsGraph parses requirements.txt and returns a flat dep list.
func buildRequirementsGraph() (map[string][]string, error) {
	deps, parseErr := parseRequirementsTxt("requirements.txt")
	if parseErr != nil {
		return nil, parseErr
	}

	graph := make(map[string][]string)
	if len(deps) > 0 {
		sort.Strings(deps)
		graph["project"] = deps
	}
	return graph, nil
}

// parseRequirementsTxt extracts package names from a requirements.txt file.
// Handles version specifiers, comments, blank lines, and -r includes.
func parseRequirementsTxt(path string) ([]string, error) {
	f, openErr := os.Open(path) //nolint:gosec // G304: path is a known manifest filename
	if openErr != nil {
		return nil, openErr
	}
	defer func() { _ = f.Close() }()

	var deps []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Skip options like -r, -e, --index-url, etc.
		if strings.HasPrefix(line, "-") {
			continue
		}
		// Extract package name before any version specifier.
		name := extractPythonPkgName(line)
		if name != "" {
			deps = append(deps, name)
		}
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return nil, scanErr
	}
	return deps, nil
}

// extractPythonPkgName extracts the package name from a requirements line.
// Handles: package==1.0, package>=1.0, package[extra]>=1.0, package ; markers
func extractPythonPkgName(line string) string {
	// Strip inline comments.
	if idx := strings.Index(line, " #"); idx >= 0 {
		line = line[:idx]
	}
	// Strip environment markers.
	if idx := strings.Index(line, ";"); idx >= 0 {
		line = line[:idx]
	}
	line = strings.TrimSpace(line)

	// Find first version specifier character.
	for i, ch := range line {
		switch ch {
		case '>', '<', '=', '!', '~', '[':
			name := strings.TrimSpace(line[:i])
			return strings.ToLower(name)
		}
	}
	return strings.ToLower(strings.TrimSpace(line))
}

// buildPyprojectGraph parses pyproject.toml for dependencies.
// Uses a simple line-based parser — no TOML library needed for this subset.
func buildPyprojectGraph(includeDevDeps bool) (map[string][]string, error) {
	data, readErr := os.ReadFile("pyproject.toml")
	if readErr != nil {
		return nil, readErr
	}

	content := string(data)
	deps := parsePyprojectDeps(content, "dependencies")

	if includeDevDeps {
		devDeps := parsePyprojectDeps(content, "dev-dependencies")
		devDeps = append(devDeps, parsePyprojectDeps(content, "dev")...)
		deps = append(deps, devDeps...)
	}

	// Deduplicate.
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
		graph["project"] = unique
	}
	return graph, nil
}

// parsePyprojectDeps extracts dependency names from a TOML array section.
// Looks for [project.dependencies], [tool.poetry.dependencies], etc.
func parsePyprojectDeps(content string, sectionSuffix string) []string {
	lines := strings.Split(content, config.NewlineLF)
	var deps []string
	inSection := false
	inArray := false

	targets := []string{
		"[project." + sectionSuffix + "]",
		"[tool.poetry." + sectionSuffix + "]",
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for section headers.
		if strings.HasPrefix(trimmed, "[") {
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
			// Also check for inline array: dependencies = [...]
			if !inArray {
				for _, key := range []string{sectionSuffix + " = [", sectionSuffix + "= [", sectionSuffix + " =[", sectionSuffix + "=["} {
					if strings.Contains(trimmed, key) {
						inArray = true
						// Parse items on this line after the opening bracket.
						idx := strings.Index(trimmed, "[")
						rest := trimmed[idx+1:]
						deps = append(deps, parsePyprojectArrayItems(rest)...)
						if strings.Contains(rest, "]") {
							inArray = false
						}
						break
					}
				}
			} else {
				deps = append(deps, parsePyprojectArrayItems(trimmed)...)
				if strings.Contains(trimmed, "]") {
					inArray = false
				}
			}
			continue
		}

		// Inside a section — look for key = "version" (Poetry style).
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Poetry style: package-name = "^1.0"
		if idx := strings.Index(trimmed, "="); idx > 0 {
			name := strings.TrimSpace(trimmed[:idx])
			// Skip python itself and metadata keys.
			lower := strings.ToLower(name)
			if lower == pythonEcosystem || lower == "name" || lower == "version" || lower == "description" {
				continue
			}
			deps = append(deps, strings.ToLower(name))
		}
	}
	return deps
}

// parsePyprojectArrayItems extracts package names from a TOML array line.
// Example: "requests>=2.0", "flask",
func parsePyprojectArrayItems(line string) []string {
	// Strip closing bracket.
	line = strings.ReplaceAll(line, "]", "")
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	var deps []string
	for _, item := range strings.Split(line, ",") {
		item = strings.TrimSpace(item)
		// Strip quotes.
		item = strings.Trim(item, `"'`)
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		name := extractPythonPkgName(item)
		if name != "" {
			deps = append(deps, name)
		}
	}
	return deps
}
