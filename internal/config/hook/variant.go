//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

// Hook variant constants: template selectors passed to Load and
// NewTemplateRef to choose the appropriate message for each trigger type.
const (
	// VariantMidSudo selects the mid-command sudo block message.
	VariantMidSudo = "mid-sudo"
	// VariantMidGitPush selects the mid-command git push block message.
	VariantMidGitPush = "mid-git-push"
	// VariantCpToBin selects the cp/mv to bin block message.
	VariantCpToBin = "cp-to-bin"
	// VariantInstallToLocalBin selects the install to ~/.local/bin block message.
	VariantInstallToLocalBin = "install-to-local-bin"
	// VariantDotSlash selects the relative path (./ctx) block message.
	VariantDotSlash = "dot-slash"
	// VariantGoRun selects the go run block message.
	VariantGoRun = "go-run"
	// VariantAbsolutePath selects the absolute path block message.
	VariantAbsolutePath = "absolute-path"
	// VariantBoth selects the template for both ceremonies missing.
	VariantBoth = "both"
	// VariantRemember selects the template for missing /ctx-remember.
	VariantRemember = "remember"
	// VariantWrapup selects the template for missing /ctx-wrap-up.
	VariantWrapup = "wrapup"
	// VariantUnimported selects the unimported journal entries variant.
	VariantUnimported = "unimported"
	// VariantUnenriched selects the unenriched journal entries variant.
	VariantUnenriched = "unenriched"
	// VariantWarning selects the generic warning variant.
	VariantWarning = "warning"
	// VariantAlert selects the alert variant.
	VariantAlert = "alert"
	// VariantBilling selects the billing threshold variant.
	VariantBilling = "billing"
	// VariantCheckpoint selects the checkpoint variant.
	VariantCheckpoint = "checkpoint"
	// VariantGate selects the gate variant.
	VariantGate = "gate"
	// VariantKeyRotation selects the key rotation variant.
	VariantKeyRotation = "key-rotation"
	// VariantMismatch selects the version mismatch variant.
	VariantMismatch = "mismatch"
	// VariantNudge selects the generic nudge variant.
	VariantNudge = "nudge"
	// VariantOversize selects the oversize threshold variant.
	VariantOversize = "oversize"
	// VariantPulse selects the heartbeat pulse variant.
	VariantPulse = "pulse"
	// VariantReminders selects the reminders variant.
	VariantReminders = "reminders"
	// VariantStale selects the staleness variant.
	VariantStale = "stale"
	// VariantWindow selects the context window variant.
	VariantWindow = "window"
)
