//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package closeout

import (
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/config/token"
	errCloseout "github.com/ActiveMemory/ctx/internal/err/closeout"
)

// renderMarkdown composes the frontmatter YAML block + body
// into the final closeout file content.
//
// Parameters:
//   - fm: parsed frontmatter for the closeout.
//   - body: rendered body markdown (everything that follows the
//     closing `---` delimiter).
//
// Returns:
//   - string: full file content with frontmatter delimiters.
//   - error: non-nil on YAML marshal failure.
func renderMarkdown(fm Frontmatter, body string) (string, error) {
	var sb strings.Builder
	sb.WriteString(frontmatterDelim)
	sb.WriteString(token.NewlineLF)
	enc, yamlErr := yaml.Marshal(fm)
	if yamlErr != nil {
		return "", errCloseout.MarshalFrontmatter(yamlErr)
	}
	sb.Write(enc)
	sb.WriteString(frontmatterDelim)
	sb.WriteString(token.DoubleNewline)
	sb.WriteString(strings.TrimRight(body, token.NewlineLF))
	sb.WriteString(token.NewlineLF)
	return sb.String(), nil
}
