//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings for system subcommands.
//
// The ctx system namespace hosts hook plumbing only. User-facing
// maintenance commands (backup, bootstrap, event, message, prune,
// resource, stats) have been promoted to top-level commands; their
// Use constants live in their own per-command files in this package.
const (
	// UseSystemBlockDangerousCommand is the cobra Use string for the system
	// block dangerous command command.
	UseSystemBlockDangerousCommand = "block-dangerous-command"
	// UseSystemBlockNonPathCtx is the cobra Use string for the system block non
	// path ctx command.
	UseSystemBlockNonPathCtx = "block-non-path-ctx"
	// UseSystemCheckBackupAge is the cobra Use string for the system check backup
	// age command.
	UseSystemCheckBackupAge = "check-backup-age"
	// UseSystemCheckCeremony is the cobra Use string for the system check
	// ceremony command.
	UseSystemCheckCeremony = "check-ceremony"
	// UseSystemCheckContextSize is the cobra Use string for the system check
	// context size command.
	UseSystemCheckContextSize = "check-context-size"
	// UseSystemCheckFreshness is the cobra Use string for the system check
	// freshness command.
	UseSystemCheckFreshness = "check-freshness"
	// UseSystemCheckHubSync is the cobra Use string for the system check hub
	// sync command.
	UseSystemCheckHubSync = "check-hub-sync"
	// UseSystemCheckJournal is the cobra Use string for the system check journal
	// command.
	UseSystemCheckJournal = "check-journal"
	// UseSystemCheckKnowledge is the cobra Use string for the system check
	// knowledge command.
	UseSystemCheckKnowledge = "check-knowledge"
	// UseSystemCheckMapStaleness is the cobra Use string for the system check map
	// staleness command.
	UseSystemCheckMapStaleness = "check-map-staleness"
	// UseSystemCheckMemoryDrift is the cobra Use string for the system check
	// memory drift command.
	UseSystemCheckMemoryDrift = "check-memory-drift"
	// UseSystemCheckPersistence is the cobra Use string for the system check
	// persistence command.
	UseSystemCheckPersistence = "check-persistence"
	// UseSystemCheckSkillDiscovery is the cobra Use string for the system check
	// skill discovery command.
	UseSystemCheckSkillDiscovery = "check-skill-discovery"
	// UseSystemCheckReminder is the cobra Use string for the system check
	// reminder command.
	UseSystemCheckReminder = "check-reminder"
	// UseSystemCheckResource is the cobra Use string for the system check
	// resource command.
	UseSystemCheckResource = "check-resource"
	// UseSystemCheckTaskCompletion is the cobra Use string for the system check
	// task completion command.
	UseSystemCheckTaskCompletion = "check-task-completion"
	// UseSystemCheckVersion is the cobra Use string for the system check version
	// command.
	UseSystemCheckVersion = "check-version"
	// UseSystemContextLoadGate is the cobra Use string for the system context
	// load gate command.
	UseSystemContextLoadGate = "context-load-gate"
	// UseSystemHeartbeat is the cobra Use string for the system heartbeat command.
	UseSystemHeartbeat = "heartbeat"
	// UseSystemMarkJournal is the cobra Use string for the system mark journal
	// command.
	UseSystemMarkJournal = "mark-journal <filename> <stage>"
	// UseSystemMarkWrappedUp is the cobra Use string for the system mark wrapped
	// up command.
	UseSystemMarkWrappedUp = "mark-wrapped-up"
	// UseSystemPause is the cobra Use string for the system pause command.
	UseSystemPause = "pause"
	// UseSystemPostCommit is the cobra Use string for the system post commit
	// command.
	UseSystemPostCommit = "post-commit"
	// UseSystemQaReminder is the cobra Use string for the system qa reminder
	// command.
	UseSystemQaReminder = "qa-reminder"
	// UseSystemResume is the cobra Use string for the system resume command.
	UseSystemResume = "resume"
	// UseSystemSessionEvent is the cobra Use string for the system session event
	// command.
	UseSystemSessionEvent = "session-event"
	// UseSystemSpecsNudge is the cobra Use string for the system specs nudge
	// command.
	UseSystemSpecsNudge = "specs-nudge"
)

// DescKeys for system subcommands.
//
// The ctx system namespace hosts hook plumbing only. DescKeys for
// promoted top-level commands live in their own per-command files.
const (
	// DescKeySystem is the description key for the system command.
	DescKeySystem = "system"
	// DescKeySystemBlockDangerousCommand is the description key for the system
	// block dangerous command command.
	DescKeySystemBlockDangerousCommand = "system.blockdangerouscommand"
	// DescKeySystemBlockNonPathCtx is the description key for the system block
	// non path ctx command.
	DescKeySystemBlockNonPathCtx = "system.blocknonpathctx"
	// DescKeySystemCheckBackupAge is the description key for the system check
	// backup age command.
	DescKeySystemCheckBackupAge = "system.checkbackupage"
	// DescKeySystemCheckCeremony is the description key for the system check
	// ceremony command.
	DescKeySystemCheckCeremony = "system.checkceremony"
	// DescKeySystemCheckContextSize is the description key for the system check
	// context size command.
	DescKeySystemCheckContextSize = "system.checkcontextsize"
	// DescKeySystemCheckFreshness is the description key for the system check
	// freshness command.
	DescKeySystemCheckFreshness = "system.checkfreshness"
	// DescKeySystemCheckHubSync is the description key for the system check
	// hub sync command.
	DescKeySystemCheckHubSync = "system.checkhubsync"
	// DescKeySystemCheckJournal is the description key for the system check
	// journal command.
	DescKeySystemCheckJournal = "system.checkjournal"
	// DescKeySystemCheckKnowledge is the description key for the system check
	// knowledge command.
	DescKeySystemCheckKnowledge = "system.checkknowledge"
	// DescKeySystemCheckMapStaleness is the description key for the system check
	// map staleness command.
	DescKeySystemCheckMapStaleness = "system.checkmapstaleness"
	// DescKeySystemCheckMemoryDrift is the description key for the system check
	// memory drift command.
	DescKeySystemCheckMemoryDrift = "system.checkmemorydrift"
	// DescKeySystemCheckPersistence is the description key for the system check
	// persistence command.
	DescKeySystemCheckPersistence = "system.checkpersistence"
	// DescKeySystemCheckSkillDiscovery is the description key for the system
	// check skill discovery command.
	DescKeySystemCheckSkillDiscovery = "system.checkskilldiscovery"
	// DescKeySystemCheckReminder is the description key for the system check
	// reminder command.
	DescKeySystemCheckReminder = "system.checkreminder"
	// DescKeySystemCheckResource is the description key for the system check
	// resource command.
	DescKeySystemCheckResource = "system.checkresource"
	// DescKeySystemCheckTaskCompletion is the description key for the system
	// check task completion command.
	DescKeySystemCheckTaskCompletion = "system.checktaskcompletion"
	// DescKeySystemCheckVersion is the description key for the system check
	// version command.
	DescKeySystemCheckVersion = "system.checkversion"
	// DescKeySystemContextLoadGate is the description key for the system context
	// load gate command.
	DescKeySystemContextLoadGate = "system.contextloadgate"
	// DescKeySystemHeartbeat is the description key for the system heartbeat
	// command.
	DescKeySystemHeartbeat = "system.heartbeat"
	// DescKeySystemMarkJournal is the description key for the system mark journal
	// command.
	DescKeySystemMarkJournal = "system.markjournal"
	// DescKeySystemMarkWrappedUp is the description key for the system mark
	// wrapped up command.
	DescKeySystemMarkWrappedUp = "system.markwrappedup"
	// DescKeySystemPause is the description key for the system pause command.
	DescKeySystemPause = "system.pause"
	// DescKeySystemPostCommit is the description key for the system post commit
	// command.
	DescKeySystemPostCommit = "system.postcommit"
	// DescKeySystemQaReminder is the description key for the system qa reminder
	// command.
	DescKeySystemQaReminder = "system.qareminder"
	// DescKeySystemResume is the description key for the system resume command.
	DescKeySystemResume = "system.resume"
	// DescKeySystemSessionEvent is the description key for the system session
	// event command.
	DescKeySystemSessionEvent = "system.sessionevent"
	// DescKeySystemSpecsNudge is the description key for the system specs nudge
	// command.
	DescKeySystemSpecsNudge = "system.specsnudge"
)
