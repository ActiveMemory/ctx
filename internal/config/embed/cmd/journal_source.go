//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings for journal source subcommands.
const (
	// UseJournalImport is the cobra Use string for the journal import command.
	UseJournalImport = "import [session-id]"
	// UseJournalLock is the cobra Use string for the journal lock command.
	UseJournalLock = "lock <pattern>"
	// UseJournalSync is the cobra Use string for the journal sync command.
	UseJournalSync = "sync"
	// UseJournalUnlock is the cobra Use string for the journal unlock command.
	UseJournalUnlock = "unlock <pattern>"
)

// DescKeys for journal source subcommands.
const (
	// DescKeyJournalImport is the description key for the journal import command.
	DescKeyJournalImport = "journal.import"
	// DescKeyJournalLock is the description key for the journal lock command.
	DescKeyJournalLock = "journal.lock"
	// DescKeyJournalSync is the description key for the journal sync command.
	DescKeyJournalSync = "journal.sync"
	// DescKeyJournalUnlock is the description key for the journal unlock command.
	DescKeyJournalUnlock = "journal.unlock"
)
