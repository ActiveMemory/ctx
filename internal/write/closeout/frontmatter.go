//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package closeout

import (
	"strings"

	"gopkg.in/yaml.v3"

	cfgCloseout "github.com/ActiveMemory/ctx/internal/config/closeout"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errCloseout "github.com/ActiveMemory/ctx/internal/err/closeout"
)

// frontmatterDelim is the YAML frontmatter open / close
// delimiter that brackets the header in a closeout file.
const frontmatterDelim = token.Separator

// splitFrontmatter splits a closeout's raw text into parsed
// frontmatter + body.
//
// Parameters:
//   - raw: full file contents.
//
// Returns:
//   - Frontmatter: parsed values.
//   - string: body (everything after closing `---`).
//   - error: [errCloseout.ErrMissingFrontmatter] or wrapped
//     YAML errors.
func splitFrontmatter(raw string) (Frontmatter, string, error) {
	lines := strings.SplitN(raw, token.NewlineLF, 2)
	if len(lines) < 2 || strings.TrimSpace(lines[0]) != frontmatterDelim {
		return Frontmatter{}, "", errCloseout.ErrMissingFrontmatter
	}
	rest := lines[1]
	openClose := token.NewlineLF + frontmatterDelim + token.NewlineLF
	idx := strings.Index(rest, openClose)
	if idx < 0 {
		// Tolerate trailing closing delim without newline.
		idx = strings.Index(rest, token.NewlineLF+frontmatterDelim)
		if idx < 0 {
			return Frontmatter{}, "", errCloseout.ErrMissingFrontmatter
		}
	}
	header := rest[:idx]
	bodyStart := idx + len(openClose)
	if bodyStart > len(rest) {
		bodyStart = len(rest)
	}
	body := rest[bodyStart:]
	body = strings.TrimLeft(body, token.NewlineLF)

	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(header), &fm); err != nil {
		return Frontmatter{}, "", errCloseout.ParseFrontmatter(err)
	}
	return fm, body, nil
}

// requireFields validates that frontmatter has the four
// always-required fields populated.
//
// Parameters:
//   - fm: parsed frontmatter.
//
// Returns:
//   - error: wraps [errCloseout.ErrMissingFields] with the
//     ordered list of missing field names; nil when all four
//     required fields are populated.
func requireFields(fm Frontmatter) error {
	missing := []string{}
	if fm.SHA == "" {
		missing = append(missing, cfgCloseout.FieldSHA)
	}
	if fm.Branch == "" {
		missing = append(missing, cfgCloseout.FieldBranch)
	}
	if fm.Mode == "" {
		missing = append(missing, cfgCloseout.FieldMode)
	}
	if fm.GeneratedAt.IsZero() {
		missing = append(missing, cfgCloseout.FieldGeneratedAt)
	}
	if len(missing) > 0 {
		return errCloseout.MissingFields(missing)
	}
	return nil
}
