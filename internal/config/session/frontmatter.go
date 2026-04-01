//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

// YAML frontmatter field names used in journal entries.
const (
	// FrontmatterTitle is the YAML frontmatter key for the entry title.
	FrontmatterTitle = "title"
	// FrontmatterDate is the YAML frontmatter key for the entry date.
	FrontmatterDate = "date"
	// FrontmatterType is the YAML frontmatter key for the session type.
	FrontmatterType = "type"
	// FrontmatterOutcome is the YAML frontmatter key for the session outcome.
	FrontmatterOutcome = "outcome"
	// FrontmatterTopics is the YAML frontmatter key for the topics list.
	FrontmatterTopics = "topics"
	// FrontmatterTechnologies is the YAML frontmatter key for
	// the technologies list.
	FrontmatterTechnologies = "technologies"
	// FrontmatterKeyFiles is the YAML frontmatter key for the key files list.
	FrontmatterKeyFiles = "key_files"
	// FrontmatterLocked is the YAML frontmatter key and journal state
	// marker for locked entries.
	FrontmatterLocked = "locked"
	// FrontmatterLockedLine is the full YAML line inserted into frontmatter
	// when a journal entry is locked.
	FrontmatterLockedLine = "locked: true  # managed by ctx"
	// Unlocked is the display label for unlocked entries.
	Unlocked = "unlocked"
)

// YAML frontmatter field keys for journal export.
const (
	FmKeyTime      = "time"
	FmKeyProject   = "project"
	FmKeyBranch    = "branch"
	FmKeyModel     = "model"
	FmKeyTokensIn  = "tokens_in"
	FmKeyTokensOut = "tokens_out"
	FmKeyID        = "session_id"
)
