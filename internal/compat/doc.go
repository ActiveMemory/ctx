//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package compat contains backward-compatibility
// integration tests that verify the hooks-and-steering
// extensions do not break existing ctx workflows when
// the new directories are absent.
//
// # What It Tests
//
// Tests exercise core commands in a clean project
// environment to confirm graceful degradation:
//
//   - [AssemblePacket] with nil steering and empty
//     skill produces a packet without Steering or
//     Skill sections, identical to pre-extension
//     behavior.
//   - Trigger discovery returns an empty map when
//     .context/hooks/ does not exist, without error.
//   - Steering resolution returns nil when
//     .context/steering/ is absent.
//   - Skill loading returns an empty string for
//     unknown skill names.
//
// # Why These Tests Exist
//
// The hooks and steering layers were added after the
// initial release. Projects created before that
// release have no .context/hooks/ or .context/steering/
// directories. These tests guarantee that every code
// path that touches those directories handles their
// absence gracefully rather than crashing or producing
// confusing errors.
//
// Every file in this package is a _test.go file except
// this doc.go. The package produces no binary output.
package compat
