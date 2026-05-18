//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Cobra Use strings for the handover command family.
const (
	// UseHandover is the cobra use string for the handover
	// parent command.
	UseHandover = "handover"
	// UseHandoverWrite is the cobra use string for
	// `ctx handover write <title>`.
	UseHandoverWrite = "write <title>"
)

// DescKeys for the handover command family. Values map to
// entries in the commands.yaml asset.
const (
	// DescKeyHandover is the description key for the handover
	// parent command.
	DescKeyHandover = "handover"
	// DescKeyHandoverWrite is the description key for
	// `ctx handover write`.
	DescKeyHandoverWrite = "handover.write"
)
