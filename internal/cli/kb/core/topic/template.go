//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package topic

import (
	"strings"

	cfgKbCli "github.com/ActiveMemory/ctx/internal/config/kb/cli"
)

// Substitute applies <NAME> / <SLUG> token replacement to the
// template body so the new file is human-meaningful from the
// first read.
//
// Parameters:
//   - body: raw template content.
//   - name: human-readable topic name.
//   - slug: kebab-case slug.
//
// Returns:
//   - string: substituted body.
func Substitute(body, name, slug string) string {
	body = strings.ReplaceAll(body, cfgKbCli.TokenTopicName, name)
	body = strings.ReplaceAll(body, cfgKbCli.TokenTopicSlug, slug)
	body = strings.ReplaceAll(body, cfgKbCli.TokenName, name)
	body = strings.ReplaceAll(body, cfgKbCli.TokenSlug, slug)
	return body
}
