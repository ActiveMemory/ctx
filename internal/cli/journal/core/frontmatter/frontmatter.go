//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package frontmatter

import (
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// Transform converts journal frontmatter to Obsidian format.
//
// List applied:
//   - topics -> tags (Obsidian-recognized key)
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
func Transform(content, sourcePath string) string {
	nl := token.NewlineLF
	fmOpen := len(token.Separator + nl)

	if !strings.HasPrefix(content, token.Separator+nl) {
		return content
	}

	endIdx := strings.Index(content[fmOpen:], nl+token.Separator+nl)
	if endIdx < 0 {
		return content
	}

	fmRaw := content[fmOpen : fmOpen+endIdx]
	afterFM := content[fmOpen+endIdx+len(nl+token.Separator+nl):]

	// Parse the original frontmatter into a generic map to preserve
	// unknown fields, then extract known fields for transformation.
	var raw map[string]any
	if yaml.Unmarshal([]byte(fmRaw), &raw) != nil {
		return content
	}

	// Build the Obsidian frontmatter
	ofm := Obsidian{}

	if v, ok := raw[session.FrontmatterTitle].(string); ok {
		ofm.Title = v
	}
	if v, ok := raw[session.FrontmatterDate].(string); ok {
		ofm.Date = v
	}
	if v, ok := raw[session.FrontmatterType].(string); ok {
		ofm.Type = v
	}
	if v, ok := raw[session.FrontmatterOutcome].(string); ok {
		ofm.Outcome = v
	}

	// topics -> tags
	ofm.Tags = ExtractStringSlice(raw, session.FrontmatterTopics)

	ofm.Technologies = ExtractStringSlice(raw, session.FrontmatterTechnologies)
	ofm.KeyFiles = ExtractStringSlice(raw, session.FrontmatterKeyFiles)

	// Add aliases from the title
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
	sb.WriteString(token.Separator + nl)
	sb.Write(out)
	sb.WriteString(token.Separator + nl)
	sb.WriteString(afterFM)

	return sb.String()
}

// ExtractStringSlice extracts a []string from a map value that may be
// []any (as returned by yaml.Unmarshal into map[string]any).
//
// Parameters:
//   - m: Source map
//   - key: Key to extract
//
// Returns:
//   - []string: Extracted strings, or nil if key is missing/empty
func ExtractStringSlice(m map[string]any, key string) []string {
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
