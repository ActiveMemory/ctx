//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/blocknonpathctx"
	sysBootstrap "github.com/ActiveMemory/ctx/internal/cli/system/cmd/bootstrap"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkaudit"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkceremony"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkcontextsize"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkfreshness"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkhubsync"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkjournal"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkknowledge"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkmapstaleness"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkmemorydrift"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkpersistence"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkreminder"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkresource"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkskilldiscovery"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checktaskcompletion"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkversion"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/contextloadgate"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/heartbeat"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/markjournal"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/markwrappedup"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/pause"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/postcommit"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/qareminder"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/resume"
	sessEvent "github.com/ActiveMemory/ctx/internal/cli/system/cmd/sessionevent"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/specsnudge"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx system" parent command.
//
// Hosts hidden Claude Code hook plumbing and agent-only commands.
// User-facing maintenance commands (prune, sysinfo, usage) are
// top-level; hook-facing commands (event, message, notify, pause,
// resume) live under "ctx hook". Both groups are registered in
// internal/bootstrap/group.go. Bootstrap remains here as
// agent-only plumbing.
//
// Hook subcommands implement Claude Code hook logic as native Go
// binaries and are not intended for direct user invocation.
//
// Returns:
//   - *cobra.Command: Parent command with hook plumbing subcommands
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeySystem, cmd.UseSystem,
		sysBootstrap.Cmd(),
		blocknonpathctx.Cmd(),
		checkceremony.Cmd(),
		checkcontextsize.Cmd(),
		checkfreshness.Cmd(),
		checkhubsync.Cmd(),
		checkjournal.Cmd(),
		checkknowledge.Cmd(),
		checkmapstaleness.Cmd(),
		checkmemorydrift.Cmd(),
		checkpersistence.Cmd(),
		checkskilldiscovery.Cmd(),
		checkaudit.Cmd(),
		checkreminder.Cmd(),
		checkresource.Cmd(),
		checktaskcompletion.Cmd(),
		checkversion.Cmd(),
		contextloadgate.Cmd(),
		heartbeat.Cmd(),
		markjournal.Cmd(),
		markwrappedup.Cmd(),
		pause.Cmd(),
		postcommit.Cmd(),
		qareminder.Cmd(),
		resume.Cmd(),
		sessEvent.Cmd(),
		specsnudge.Cmd(),
	)
}
