//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings for pad subcommands.
const (
	// UsePadAdd is the cobra Use string for the pad add command.
	UsePadAdd = "add TEXT"
	// UsePadEdit is the cobra Use string for the pad edit command.
	UsePadEdit = "edit N [TEXT]"
	// UsePadExport is the cobra Use string for the pad export command.
	UsePadExport = "export [DIR]"
	// UsePadImport is the cobra Use string for the pad import command.
	UsePadImport = "import FILE"
	// UsePadMerge is the cobra Use string for the pad merge command.
	UsePadMerge = "merge FILE..."
	// UsePadMv is the cobra Use string for the pad mv command.
	UsePadMv = "mv N M"
	// UsePadResolve is the cobra Use string for the pad resolve command.
	UsePadResolve = "resolve"
	// UsePadRm is the cobra Use string for the pad rm command.
	UsePadRm = "rm N"
	// UsePadShow is the cobra Use string for the pad show command.
	UsePadShow = "show N"
)

// DescKeys for pad subcommands.
const (
	// DescKeyPad is the description key for the pad command.
	DescKeyPad = "pad"
	// DescKeyPadAdd is the description key for the pad add command.
	DescKeyPadAdd = "pad.add"
	// DescKeyPadEdit is the description key for the pad edit command.
	DescKeyPadEdit = "pad.edit"
	// DescKeyPadExport is the description key for the pad export command.
	DescKeyPadExport = "pad.export"
	// DescKeyPadImp is the description key for the pad imp command.
	DescKeyPadImp = "pad.root"
	// DescKeyPadMerge is the description key for the pad merge command.
	DescKeyPadMerge = "pad.merge"
	// DescKeyPadMv is the description key for the pad mv command.
	DescKeyPadMv = "pad.mv"
	// DescKeyPadResolve is the description key for the pad resolve command.
	DescKeyPadResolve = "pad.resolve"
	// DescKeyPadRm is the description key for the pad rm command.
	DescKeyPadRm = "pad.rm"
	// DescKeyPadShow is the description key for the pad show command.
	DescKeyPadShow = "pad.show"
)
