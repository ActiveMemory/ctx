//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package flag holds the **lookup keys** for every CLI flag's
// help text, the short blurb cobra prints next to `--name`
// in `ctx <command> --help` output.
//
// The package is the flag-help half of the same two-step
// indirection used by [internal/config/embed/text] and
// [internal/config/embed/cmd]:
//
//  1. **Here**: `DescKeyXxx` Go constants, one per flag
//     across every command in the binary.
//  2. **In** [internal/assets/commands/text/*.yaml]: the
//     actual help string. Resolved at run-time via
//     [internal/assets/read/desc.Flag](key).
//
// The split keeps flag wording editable without a Go
// rebuild, lets the audit suite catch typos at CI time, and
// makes per-locale flag help structurally possible.
//
// # File Layout: One Command per File
//
// Each file groups the flag-key constants for one command
// (`add.go`, `agent.go`, `backup.go`, …). Within a file,
// constants follow the alphabetical-by-flag-name order
// cobra itself uses.
//
// # Naming Convention
//
// `DescKey<Command><FlagName>` for the constant; the YAML
// key is the dotted form `<command>.<flagname>`. The audit
// suite (`desckey_namespace_test`) verifies every constant
// has a matching YAML entry and every YAML entry has a
// matching constant.
//
// # Usage
//
//	import (
//	    "github.com/ActiveMemory/ctx/internal/assets/read/desc"
//	    "github.com/ActiveMemory/ctx/internal/config/embed/flag"
//	)
//	c.Flags().Bool("dry-run", false, desc.Flag(flag.DescKeyAddDryRun))
//
// In practice, most flag binding goes through
// [internal/flagbind] which already knows how to look up
// the desc key, so callers rarely call `desc.Flag` directly.
package flag
