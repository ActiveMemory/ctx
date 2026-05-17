//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for kb source-coverage ledger error wrappers.
const (
	// DescKeyErrKbSourcecoverageIllegalTransition wraps the
	// illegal-transition sentinel with from/to/source operands.
	DescKeyErrKbSourcecoverageIllegalTransition = "err.kb.sourcecoverage.illegal-transition"
	// DescKeyErrKbSourcecoverageUnknownSource wraps the
	// unknown-source sentinel with source/state operands.
	DescKeyErrKbSourcecoverageUnknownSource = "err.kb.sourcecoverage.unknown-source"
	// DescKeyErrKbSourcecoverageReadLedger wraps a ledger-read
	// failure.
	DescKeyErrKbSourcecoverageReadLedger = "err.kb.sourcecoverage.read-ledger"
	// DescKeyErrKbSourcecoverageMkdirLedgerDir wraps a parent-dir
	// mkdir failure for the ledger.
	DescKeyErrKbSourcecoverageMkdirLedgerDir = "err.kb.sourcecoverage.mkdir-ledger-dir"
	// DescKeyErrKbSourcecoverageWriteLedger wraps a write
	// failure on the ledger.
	DescKeyErrKbSourcecoverageWriteLedger = "err.kb.sourcecoverage.write-ledger"
)
