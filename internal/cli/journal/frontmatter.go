//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/config"
)

// obsidianFrontmatter represents the YAML frontmatter for Obsidian vault
// entries. Extends journalFrontmatter with Obsidian-specific fields.
type obsidianFrontmatter struct {
	Title        string   `yaml:"title"`
	Date         string   `yaml:"date"`
	Type         string   `yaml:"type,omitempty"`
	Outcome      string   `yaml:"outcome,omitempty"`
	Tags         []string `yaml:"tags,omitempty"`
	Technologies []string `yaml:"technologies,omitempty"`
	KeyFiles     []string `yaml:"key_files,omitempty"`
	Aliases      []string `yaml:"aliases,omitempty"`
	SourceFile   string   `yaml:"source_file,omitempty"`
}

// transformFrontmatter converts journal frontmatter to Obsidian format.
//
// Changes applied:
//   - topics → tags (Obsidian-recognized key)
//   - aliases added from title (makes entries findable by name)
//   - source_file added with the relative path to the source entry
//   - technologies preserved as custom property
//
// Parameters:
//   - content: Full Markdown content with YAML frontmatter
//   - sourcePath: Relative path to the source journal file
//
// Returns:
//   - string: Content with transformed frontmatter
func transformFrontmatter(content, sourcePath string) string {
	nl := config.NewlineLF
	fmOpen := len(config.Separator + nl)

	if !strings.HasPrefix(content, config.Separator+nl) {
		return content
	}

	endIdx := strings.Index(content[fmOpen:], nl+config.Separator+nl)
	if endIdx < 0 {
		return content
	}

	fmRaw := content[fmOpen : fmOpen+endIdx]
	afterFM := content[fmOpen+endIdx+len(nl+config.Separator+nl):]

	// Parse the original frontmatter into a generic map to preserve
	// unknown fields, then extract known fields for transformation.
	var raw map[string]any
	if yaml.Unmarshal([]byte(fmRaw), &raw) != nil {
		return content
	}

	// Build the Obsidian frontmatter
	ofm := obsidianFrontmatter{}

	if v, ok := raw["title"].(string); ok {
		ofm.Title = v
	}
	if v, ok := raw["date"].(string); ok {
		ofm.Date = v
	}
	if v, ok := raw["type"].(string); ok {
		ofm.Type = v
	}
	if v, ok := raw["outcome"].(string); ok {
		ofm.Outcome = v
	}

	// topics → tags
	ofm.Tags = extractStringSlice(raw, "topics")

	ofm.Technologies = extractStringSlice(raw, "technologies")
	ofm.KeyFiles = extractStringSlice(raw, "key_files")

	// Add aliases from title
	if ofm.Title != "" {
		ofm.Aliases = []string{ofm.Title}
	}

	// Add source file reference
	if sourcePath != "" {
		ofm.SourceFile = sourcePath
	}

	out, marshalErr := yaml.Marshal(&ofm)
	if marshalErr != nil {
		return content
	}

	var sb strings.Builder
	sb.WriteString(config.Separator + nl)
	sb.Write(out)
	sb.WriteString(config.Separator + nl)
	sb.WriteString(afterFM)

	return sb.String()
}

// extractStringSlice extracts a []string from a map value that may be
// []any (as returned by yaml.Unmarshal into map[string]any).
//
// Parameters:
//   - m: Source map
//   - key: Key to extract
//
// Returns:
//   - []string: Extracted strings, or nil if key is missing/empty
func extractStringSlice(m map[string]any, key string) []string {
	val, ok := m[key]
	if !ok {
		return nil
	}

	items, ok := val.([]any)
	if !ok {
		return nil
	}

	result := make([]string, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}

	if len(result) == 0 {
		return nil
	}
	return result
}
