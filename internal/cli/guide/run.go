//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package guide

import (
	"fmt"

	"github.com/spf13/cobra"
)

const defaultGuide = `ctx — persistent AI context

GETTING STARTED
  ctx init              Create .context/ directory with templates
  ctx status            Show context health summary
  ctx doctor            Diagnose configuration issues

TRACKING DECISIONS & KNOWLEDGE
  ctx add -t TYPE       Add a decision, learning, convention, or task
  ctx complete          Mark a task as done in TASKS.md
  ctx decisions reindex Rebuild the DECISIONS.md index table

BROWSING HISTORY
  ctx recall list       List exported session transcripts
  ctx recall show ID    Read a specific session transcript
  ctx journal site      Generate a browsable journal site

AI CONTEXT
  ctx agent             Load AI-optimized context packet (use --budget)
  ctx load              Output raw assembled context
  ctx drift             Detect stale or invalid context

MAINTENANCE
  ctx compact           Archive completed tasks, trim context
  ctx sync              Reconcile codebase changes with docs
  ctx pad               Encrypted scratchpad for sensitive notes

KEY SKILLS
  /ctx-commit           Commit with context capture
  /ctx-implement        Execute a plan step-by-step
  /ctx-next             Suggest what to work on next
  /ctx-reflect          Surface persist-worthy items
  /ctx-wrap-up          End-of-session persistence ceremony
  /ctx-drift            Detect and fix context drift
  /ctx-remember         Recall project context
  /ctx-recall           Browse session history

RECIPES
  Start session:      ctx agent --budget 4000
  Record decision:    ctx add -t decision "Use PostgreSQL for persistence"
  Check health:       ctx status && ctx drift
  End session:        invoke /ctx-wrap-up

Full listings: ctx guide --skills | ctx guide --commands
`

// runGuide dispatches to the appropriate output based on flags.
func runGuide(cmd *cobra.Command, showSkills, showCommands bool) error {
	switch {
	case showSkills:
		return listSkills(cmd)
	case showCommands:
		return listCommands(cmd)
	default:
		cmd.Print(defaultGuide)
		return nil
	}
}

// listCommands prints all non-hidden subcommands from the root.
func listCommands(cmd *cobra.Command) error {
	root := cmd.Root()
	cmd.Println("CLI Commands:")
	cmd.Println()
	for _, c := range root.Commands() {
		if c.Hidden {
			continue
		}
		cmd.Println(fmt.Sprintf("  %-14s %s", c.Name(), c.Short))
	}
	return nil
}
