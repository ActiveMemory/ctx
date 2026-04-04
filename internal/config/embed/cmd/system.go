//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings for system subcommands.
const (
	// UseSystemBackup is the cobra Use string for the system backup command.
	UseSystemBackup = "backup"
	// UseSystemBlockDangerousCommands is the cobra Use string for the system
	// block dangerous commands command.
	UseSystemBlockDangerousCommands = "block-dangerous-commands"
	// UseSystemBlockNonPathCtx is the cobra Use string for the system block non
	// path ctx command.
	UseSystemBlockNonPathCtx = "block-non-path-ctx"
	// UseSystemBootstrap is the cobra Use string for the system bootstrap command.
	UseSystemBootstrap = "bootstrap"
	// UseSystemCheckBackupAge is the cobra Use string for the system check backup
	// age command.
	UseSystemCheckBackupAge = "check-backup-age"
	// UseSystemCheckCeremonies is the cobra Use string for the system check
	// ceremonies command.
	UseSystemCheckCeremonies = "check-ceremonies"
	// UseSystemCheckContextSize is the cobra Use string for the system check
	// context size command.
	UseSystemCheckContextSize = "check-context-size"
	// UseSystemCheckFreshness is the cobra Use string for the system check
	// freshness command.
	UseSystemCheckFreshness = "check-freshness"
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
	// UseSystemCheckReminders is the cobra Use string for the system check
	// reminders command.
	UseSystemCheckReminders = "check-reminders"
	// UseSystemCheckResources is the cobra Use string for the system check
	// resources command.
	UseSystemCheckResources = "check-resources"
	// UseSystemCheckTaskCompletion is the cobra Use string for the system check
	// task completion command.
	UseSystemCheckTaskCompletion = "check-task-completion"
	// UseSystemCheckVersion is the cobra Use string for the system check version
	// command.
	UseSystemCheckVersion = "check-version"
	// UseSystemContextLoadGate is the cobra Use string for the system context
	// load gate command.
	UseSystemContextLoadGate = "context-load-gate"
	// UseSystemEvents is the cobra Use string for the system events command.
	UseSystemEvents = "events"
	// UseSystemHeartbeat is the cobra Use string for the system heartbeat command.
	UseSystemHeartbeat = "heartbeat"
	// UseSystemMarkJournal is the cobra Use string for the system mark journal
	// command.
	UseSystemMarkJournal = "mark-journal <filename> <stage>"
	// UseSystemMarkWrappedUp is the cobra Use string for the system mark wrapped
	// up command.
	UseSystemMarkWrappedUp = "mark-wrapped-up"
	// UseSystemMessage is the cobra Use string for the system message command.
	UseSystemMessage = "message"
	// UseSystemMessageEdit is the cobra Use string for the system message edit
	// command.
	UseSystemMessageEdit = "edit <hook> <variant>"
	// UseSystemMessageList is the cobra Use string for the system message list
	// command.
	UseSystemMessageList = "list"
	// UseSystemMessageReset is the cobra Use string for the system message reset
	// command.
	UseSystemMessageReset = "reset <hook> <variant>"
	// UseSystemMessageShow is the cobra Use string for the system message show
	// command.
	UseSystemMessageShow = "show <hook> <variant>"
	// UseSystemPause is the cobra Use string for the system pause command.
	UseSystemPause = "pause"
	// UseSystemPostCommit is the cobra Use string for the system post commit
	// command.
	UseSystemPostCommit = "post-commit"
	// UseSystemPrune is the cobra Use string for the system prune command.
	UseSystemPrune = "prune"
	// UseSystemQaReminder is the cobra Use string for the system qa reminder
	// command.
	UseSystemQaReminder = "qa-reminder"
	// UseSystemResources is the cobra Use string for the system resources command.
	UseSystemResources = "resources"
	// UseSystemResume is the cobra Use string for the system resume command.
	UseSystemResume = "resume"
	// UseSystemSessionEvent is the cobra Use string for the system session event
	// command.
	UseSystemSessionEvent = "session-event"
	// UseSystemSpecsNudge is the cobra Use string for the system specs nudge
	// command.
	UseSystemSpecsNudge = "specs-nudge"
	// UseSystemStats is the cobra Use string for the system stats command.
	UseSystemStats = "stats"
)

// DescKeys for system subcommands.
const (
	// DescKeySystem is the description key for the system command.
	DescKeySystem = "system"
	// DescKeySystemBackup is the description key for the system backup command.
	DescKeySystemBackup = "system.backup"
	// DescKeySystemBlockDangerousCommands is the description key for the system
	// block dangerous commands command.
	DescKeySystemBlockDangerousCommands = "system.blockdangerouscommands"
	// DescKeySystemBlockNonPathCtx is the description key for the system block
	// non path ctx command.
	DescKeySystemBlockNonPathCtx = "system.blocknonpathctx"
	// DescKeySystemBootstrap is the description key for the system bootstrap
	// command.
	DescKeySystemBootstrap = "system.bootstrap"
	// DescKeySystemCheckBackupAge is the description key for the system check
	// backup age command.
	DescKeySystemCheckBackupAge = "system.checkbackupage"
	// DescKeySystemCheckCeremonies is the description key for the system check
	// ceremonies command.
	DescKeySystemCheckCeremonies = "system.checkceremonies"
	// DescKeySystemCheckContextSize is the description key for the system check
	// context size command.
	DescKeySystemCheckContextSize = "system.checkcontextsize"
	// DescKeySystemCheckFreshness is the description key for the system check
	// freshness command.
	DescKeySystemCheckFreshness = "system.checkfreshness"
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
	// DescKeySystemCheckReminders is the description key for the system check
	// reminders command.
	DescKeySystemCheckReminders = "system.checkreminders"
	// DescKeySystemCheckResources is the description key for the system check
	// resources command.
	DescKeySystemCheckResources = "system.checkresources"
	// DescKeySystemCheckTaskCompletion is the description key for the system
	// check task completion command.
	DescKeySystemCheckTaskCompletion = "system.checktaskcompletion"
	// DescKeySystemCheckVersion is the description key for the system check
	// version command.
	DescKeySystemCheckVersion = "system.checkversion"
	// DescKeySystemContextLoadGate is the description key for the system context
	// load gate command.
	DescKeySystemContextLoadGate = "system.contextloadgate"
	// DescKeySystemEvents is the description key for the system events command.
	DescKeySystemEvents = "system.events"
	// DescKeySystemHeartbeat is the description key for the system heartbeat
	// command.
	DescKeySystemHeartbeat = "system.heartbeat"
	// DescKeySystemMarkJournal is the description key for the system mark journal
	// command.
	DescKeySystemMarkJournal = "system.markjournal"
	// DescKeySystemMarkWrappedUp is the description key for the system mark
	// wrapped up command.
	DescKeySystemMarkWrappedUp = "system.markwrappedup"
	// DescKeySystemMessage is the description key for the system message command.
	DescKeySystemMessage = "system.message"
	// DescKeySystemMessageEdit is the description key for the system message edit
	// command.
	DescKeySystemMessageEdit = "system.message.edit"
	// DescKeySystemMessageList is the description key for the system message list
	// command.
	DescKeySystemMessageList = "system.message.list"
	// DescKeySystemMessageReset is the description key for the system message
	// reset command.
	DescKeySystemMessageReset = "system.message.reset"
	// DescKeySystemMessageShow is the description key for the system message show
	// command.
	DescKeySystemMessageShow = "system.message.show"
	// DescKeySystemPause is the description key for the system pause command.
	DescKeySystemPause = "system.pause"
	// DescKeySystemPostCommit is the description key for the system post commit
	// command.
	DescKeySystemPostCommit = "system.postcommit"
	// DescKeySystemPrune is the description key for the system prune command.
	DescKeySystemPrune = "system.prune"
	// DescKeySystemQaReminder is the description key for the system qa reminder
	// command.
	DescKeySystemQaReminder = "system.qareminder"
	// DescKeySystemResources is the description key for the system resources
	// command.
	DescKeySystemResources = "system.resources"
	// DescKeySystemResume is the description key for the system resume command.
	DescKeySystemResume = "system.resume"
	// DescKeySystemSessionEvent is the description key for the system session
	// event command.
	DescKeySystemSessionEvent = "system.sessionevent"
	// DescKeySystemSpecsNudge is the description key for the system specs nudge
	// command.
	DescKeySystemSpecsNudge = "system.specsnudge"
	// DescKeySystemStats is the description key for the system stats command.
	DescKeySystemStats = "system.stats"
)
