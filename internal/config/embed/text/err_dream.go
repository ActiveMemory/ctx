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
)
