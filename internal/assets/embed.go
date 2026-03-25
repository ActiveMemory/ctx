//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

import (
	"embed"
)

//go:embed claude/.claude-plugin/plugin.json claude/CLAUDE.md
//go:embed claude/skills/*/references/*.md claude/skills/*/SKILL.md
//go:embed context/*.md project/* entry-templates/*.md hooks/*.md
//go:embed hooks/copilot-cli/*.json hooks/copilot-cli/*.md hooks/copilot-cli/scripts/*.sh hooks/copilot-cli/scripts/*.ps1
//go:embed hooks/messages/*/*.txt hooks/messages/registry.yaml
//go:embed prompt-templates/*.md ralph/*.md schema/*.json why/*.md
//go:embed permissions/*.txt commands/*.yaml commands/text/*.yaml journal/*.css
var FS embed.FS
