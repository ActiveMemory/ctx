//	/    Context:                     https://ctx.ist
//
// ,'`./    do you remember?
//
//	`.,'\
//	  \    Copyright 2026-present Context contributors.
//	                SPDX-License-Identifier: Apache-2.

package config

// Field name constants for structured entry attributes.
//
// These are used in validation error messages and as attribute names
// in context-update XML tags for decisions and learnings.
const (
	// FieldContext is the background/situation field for decisions and learnings.
	FieldContext = "context"
	// FieldRationale is the reasoning field for decisions (why this choice).
	FieldRationale = "rationale"
	// FieldConsequence is the outcomes field for decisions (what changes).
	FieldConsequence = "consequences"
	// FieldApplication is the usage field for learnings (how to apply going forward).
	FieldApplication = "application"
	// FieldLesson is the insight field for learnings (the key takeaway).
	FieldLesson = "lesson"
)
