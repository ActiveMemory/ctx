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
//go:embed codex/.codex-plugin/plugin.json codex/.mcp.json
//go:embed codex/hooks/hooks.json codex/skills/*/SKILL.md
//go:embed codex/skills/*/references/*
//go:embed context/*.md project/* entry-templates/*.md integrations/agents.md
//go:embed integrations/copilot/*.md
//go:embed integrations/copilot-cli/*.json integrations/copilot-cli/*.md
//go:embed integrations/copilot-cli/skills/*/SKILL.md
//go:embed integrations/opencode/plugin/index.ts
//go:embed integrations/opencode/skills/*/SKILL.md
//go:embed integrations/opencode/skills/*/references/*
//go:embed integrations/pi/extension/ctx.ts
//go:embed integrations/pi/skills/*/SKILL.md
//go:embed integrations/pi/skills/*/references/*
//go:embed hooks/messages/*/*.txt hooks/messages/registry.yaml hooks/trace/*.sh
//go:embed schema/*.json why/*.md
//go:embed permissions/*.txt commands/*.yaml commands/text/*.yaml journal/*.css
//go:embed i18n/placeholders/*.yaml
//go:embed kb/templates/ingest/*.md kb/templates/ingest/schemas/*.md
//go:embed kb/templates/kb/index.md kb/templates/kb/topics/_template/index.md
var FS embed.FS
