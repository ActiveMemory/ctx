//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package path resolves on-disk locations for the handover
// mechanism. Handover is the session-to-session glue created
// by `/ctx-wrap-up` and read by `/ctx-remember`; it is
// universal to every ctx project and does not depend on the
// editorial pipeline.
//
// # Public Surface
//
//   - [Dir]: returns the `.context/handovers/` directory.
//
// The closeout-archive directory used by the optional fold
// pass lives under
// [github.com/ActiveMemory/ctx/internal/cli/kb/core/path]
// (closeouts themselves are editorial-pipeline artifacts; the
// handover writer happens to be their consumer when KB is
// present).
package path
