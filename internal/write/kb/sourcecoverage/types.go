//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sourcecoverage

import "time"

// Row is one source-coverage ledger entry. Fields match the
// schema at
// `internal/assets/kb/templates/ingest/schemas/source-coverage.md`.
type Row struct {
	// Source is the short-name from `source-map.md` that
	// identifies the source uniquely within this kb.
	Source string
	// Topic is the kb-topic slug this source contributes to,
	// or the literal "n/a" for non-topic passes.
	Topic string
	// State is one of the state-name constants in
	// [github.com/ActiveMemory/ctx/internal/config/kb] (e.g.
	// `cfgKB.StateAdmitted`).
	State string
	// EVCoverage names the EV-### range minted from this
	// source, e.g. "EV-018..EV-034", or "none".
	EVCoverage string
	// Residue is a short free-text note describing what is
	// not-yet-covered by the page(s) backed by this source.
	Residue string
	// NextAction is the exact resumption invocation that would
	// advance this row, e.g.
	// "/ctx-kb-ingest cursor-hooks (resume topic-page)".
	NextAction string
	// Updated is the timestamp of the most recent pass that
	// touched this row.
	Updated time.Time
}

// transition keys the allowed-transitions map declared in
// transition.go. The from / to fields name the state pair as
// declared in [github.com/ActiveMemory/ctx/internal/config/kb].
type transition struct {
	from string
	to   string
}
