//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings for journal source subcommands.
const (
	UseJournalImport = "import [session-id]"
	UseJournalLock   = "lock <pattern>"
	UseJournalSync   = "sync"
	UseJournalUnlock = "unlock <pattern>"
)

// DescKeys for journal source subcommands.
const (
	DescKeyJournalImport = "journal.import"
	DescKeyJournalLock   = "journal.lock"
	DescKeyJournalSync   = "journal.sync"
	DescKeyJournalUnlock = "journal.unlock"
)
