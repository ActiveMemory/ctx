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
//go:embed context/*.md project/* entry-templates/*.md integrations/agents.md
//go:embed integrations/copilot/*.md
//go:embed integrations/copilot-cli/*.json integrations/copilot-cli/*.md
//go:embed integrations/copilot-cli/scripts/*.sh
//go:embed integrations/copilot-cli/scripts/*.ps1
//go:embed integrations/copilot-cli/skills/*/SKILL.md
//go:embed hooks/messages/*/*.txt hooks/messages/registry.yaml hooks/trace/*.sh
//go:embed schema/*.json why/*.md
//go:embed permissions/*.txt commands/*.yaml commands/text/*.yaml journal/*.css
var FS embed.FS
