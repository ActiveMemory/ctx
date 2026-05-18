//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handover

import (
	"time"

	"github.com/ActiveMemory/ctx/internal/entity"
)

// Entry is the caller-supplied content for a new handover. SHA
// and Branch are resolved automatically when empty.
type Entry struct {
	// Title is the short slug used to compose the filename.
	Title string
	// Summary records what happened this session in past tense.
	// Required and validated non-placeholder at the CLI layer.
	Summary string
	// Next records what the next agent should do FIRST. Future
	// tense, specific. Required and validated non-placeholder
	// at the CLI layer.
	Next string
	// Highlights records notable artifacts produced this
	// session. Optional.
	Highlights string
	// OpenQuestions lists things that remain undecided.
	// Optional.
	OpenQuestions string
	// CommitOverride forces a specific commit SHA into the
	// Provenance line, bypassing gitmeta.ResolveHead. Used by
	// the --commit flag for CI replay.
	CommitOverride string
	// NoFold skips closeout consumption (mid-session
	// checkpoint). Defaults to false; default is "fold".
	NoFold bool
}

// Frontmatter holds the four required frontmatter fields of a
// handover file.
type Frontmatter struct {
	SHA         string    `yaml:"sha"`
	Branch      string    `yaml:"branch"`
	GeneratedAt time.Time `yaml:"generated-at"`
	Title       string    `yaml:"title"`
}

// File pairs a handover's on-disk path with parsed frontmatter
// and the raw body bytes.
type File struct {
	Path        string
	Frontmatter Frontmatter
	Body        string
}

// Result reports what Write actually did.
type Result struct {
	// File is the newly-written handover.
	File File
	// FoldedCloseouts lists closeouts folded into the handover
	// and archived. Empty when NoFold is true or no closeouts
	// were postdated.
	FoldedCloseouts []entity.CloseoutFile
	// MalformedCloseouts lists paths of closeout files that
	// failed to parse during fold. Non-fatal; surfaced for the
	// doctor advisory.
	MalformedCloseouts []string
}
