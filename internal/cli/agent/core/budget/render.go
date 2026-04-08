//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package budget

import (
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/io"
)

// RenderMarkdownPacket renders an assembled packet as Markdown.
//
// Parameters:
//   - pkt: Assembled packet to render
//
// Returns:
//   - string: Formatted Markdown output
func RenderMarkdownPacket(pkt *AssembledPacket) string {
	var sb strings.Builder
	nl := token.NewlineLF

	sb.WriteString(desc.Text(text.DescKeyAgentPacketTitle) + nl)
	io.SafeFprintf(&sb,
		desc.Text(text.DescKeyAgentPacketMeta),
		time.Now().UTC().Format(time.RFC3339), pkt.Budget, pkt.TokensUsed)
	sb.WriteString(nl + nl)

	// Read order
	sb.WriteString(desc.Text(text.DescKeyAgentSectionReadOrder) + nl)
	for i, path := range pkt.ReadOrder {
		io.SafeFprintf(&sb,
			desc.Text(text.DescKeyWriteAgentNumberedItem), i+1, path)
		sb.WriteString(nl)
	}
	sb.WriteString(nl)

	// Constitution
	if len(pkt.Constitution) > 0 {
		sb.WriteString(desc.Text(text.DescKeyAgentSectionConstitution) + nl)
		for _, rule := range pkt.Constitution {
			io.SafeFprintf(&sb,
				desc.Text(text.DescKeyWriteAgentBulletItem), rule)
			sb.WriteString(nl)
		}
		sb.WriteString(nl)
	}

	// Tasks
	if len(pkt.Tasks) > 0 {
		sb.WriteString(desc.Text(text.DescKeyAgentSectionTasks) + nl)
		for _, t := range pkt.Tasks {
			sb.WriteString(t + nl)
		}
		sb.WriteString(nl)
	}

	// Conventions
	if len(pkt.Conventions) > 0 {
		sb.WriteString(desc.Text(text.DescKeyAgentSectionConventions) + nl)
		for _, conv := range pkt.Conventions {
			io.SafeFprintf(&sb,
				desc.Text(text.DescKeyWriteAgentBulletItem), conv)
			sb.WriteString(nl)
		}
		sb.WriteString(nl)
	}

	// Decisions (full body)
	if len(pkt.Decisions) > 0 {
		sb.WriteString(desc.Text(text.DescKeyAgentSectionDecisions) + nl)
		for _, dec := range pkt.Decisions {
			sb.WriteString(dec + nl + nl)
		}
	}

	// Learnings (full body)
	if len(pkt.Learnings) > 0 {
		sb.WriteString(desc.Text(text.DescKeyAgentSectionLearnings) + nl)
		for _, learn := range pkt.Learnings {
			sb.WriteString(learn + nl + nl)
		}
	}

	// Summaries
	if len(pkt.Summaries) > 0 {
		sb.WriteString(desc.Text(text.DescKeyAgentSectionSummaries) + nl)
		for _, s := range pkt.Summaries {
			io.SafeFprintf(&sb,
				desc.Text(text.DescKeyWriteAgentBulletItem), s)
			sb.WriteString(nl)
		}
		sb.WriteString(nl)
	}

	// Steering
	if len(pkt.Steering) > 0 {
		sb.WriteString(desc.Text(text.DescKeyAgentSectionSteering) + nl)
		for _, s := range pkt.Steering {
			sb.WriteString(s + nl + nl)
		}
	}

	// Shared hub entries
	if len(pkt.Shared) > 0 {
		sb.WriteString("## Shared Knowledge" + nl)
		for _, s := range pkt.Shared {
			sb.WriteString(s + nl + nl)
		}
	}

	// Skill
	if pkt.Skill != "" {
		sb.WriteString(desc.Text(text.DescKeyAgentSectionSkill) + nl)
		sb.WriteString(pkt.Skill + nl + nl)
	}

	sb.WriteString(pkt.Instruction + nl)

	return sb.String()
}
