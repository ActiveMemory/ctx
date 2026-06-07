//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for dream engine errors.
const (
	// DescKeyErrDreamCheckIgnore is the text key for err dream
	// check-ignore exec failure messages.
	DescKeyErrDreamCheckIgnore = "err.dream.check-ignore"
	// DescKeyErrDreamWriteScope is the text key for err dream
	// write-scope guard refusal messages.
	DescKeyErrDreamWriteScope = "err.dream.write-scope"
	// DescKeyErrDreamLeak is the text key for err dream don't-leak
	// guard refusal (tracked path) messages.
	DescKeyErrDreamLeak = "err.dream.leak"
	// DescKeyErrDreamResolveRoot is the text key for err dream
	// project-root resolution failure messages.
	DescKeyErrDreamResolveRoot = "err.dream.resolve-root"
	// DescKeyErrDreamRelPath is the text key for err dream relative
	// path computation failure messages.
	DescKeyErrDreamRelPath = "err.dream.rel-path"
	// DescKeyErrDreamReadState is the text key for err dream state
	// file read failure messages.
	DescKeyErrDreamReadState = "err.dream.read-state"
	// DescKeyErrDreamWriteState is the text key for err dream state
	// file write failure messages.
	DescKeyErrDreamWriteState = "err.dream.write-state"
	// DescKeyErrDreamMarshalState is the text key for err dream state
	// JSON marshal failure messages.
	DescKeyErrDreamMarshalState = "err.dream.marshal-state"
	// DescKeyErrDreamUnmarshalState is the text key for err dream
	// state JSON unmarshal failure messages.
	DescKeyErrDreamUnmarshalState = "err.dream.unmarshal-state"
	// DescKeyErrDreamAppendLedger is the text key for err dream
	// ledger append failure messages.
	DescKeyErrDreamAppendLedger = "err.dream.append-ledger"
	// DescKeyErrDreamReadLedger is the text key for err dream ledger
	// read failure messages.
	DescKeyErrDreamReadLedger = "err.dream.read-ledger"
	// DescKeyErrDreamMarshalEntry is the text key for err dream
	// ledger entry JSON marshal failure messages.
	DescKeyErrDreamMarshalEntry = "err.dream.marshal-entry"
	// DescKeyErrDreamMkdir is the text key for err dream notebook
	// directory creation failure messages.
	DescKeyErrDreamMkdir = "err.dream.mkdir"
	// DescKeyErrDreamInvalidProposal is the text key for err dream
	// invalid-proposal validation messages.
	DescKeyErrDreamInvalidProposal = "err.dream.invalid-proposal"
	// DescKeyErrDreamBackupFailed is the text key for err dream
	// backup-before-mutate failure messages.
	DescKeyErrDreamBackupFailed = "err.dream.backup-failed"
	// DescKeyErrDreamExecutorNotFound is the text key for err dream
	// executor-not-on-PATH fail-loud messages.
	DescKeyErrDreamExecutorNotFound = "err.dream.executor-not-found"
	// DescKeyErrDreamExecutorRun is the text key for err dream
	// executor-run-failure fail-loud messages.
	DescKeyErrDreamExecutorRun = "err.dream.executor-run"
	// DescKeyErrDreamGuardRefused is the text key for err dream
	// guard-refusal messages (the registry-sourced reason verbatim).
	DescKeyErrDreamGuardRefused = "err.dream.guard-refused"
	// DescKeyErrDreamLockAcquire is the text key for err dream lock
	// acquisition failure messages.
	DescKeyErrDreamLockAcquire = "err.dream.lock-acquire"
	// DescKeyErrDreamMoveSource is the text key for err dream
	// source-relocation failure messages.
	DescKeyErrDreamMoveSource = "err.dream.move-source"
	// DescKeyErrDreamProposalNotFound is the text key for err dream
	// proposal-not-found messages.
	DescKeyErrDreamProposalNotFound = "err.dream.proposal-not-found"
	// DescKeyErrDreamReadProposals is the text key for err dream
	// proposals read failure messages.
	DescKeyErrDreamReadProposals = "err.dream.read-proposals"
	// DescKeyErrDreamReadSource is the text key for err dream source
	// read failure messages.
	DescKeyErrDreamReadSource = "err.dream.read-source"
	// DescKeyErrDreamScanIdeas is the text key for err dream
	// ideas-scan failure messages.
	DescKeyErrDreamScanIdeas = "err.dream.scan-ideas"
	// DescKeyErrDreamUnknownAction is the text key for err dream
	// unknown-action messages.
	DescKeyErrDreamUnknownAction = "err.dream.unknown-action"
)
