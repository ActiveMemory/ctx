//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package file

// Common filenames.
const (
	// Readme is the standard README filename.
	Readme = "README.md"
	// Index is the standard index filename for generated sites.
	Index = "index.md"
	// SchemaDrift is the schema drift report in .context/reports/.
	SchemaDrift = "schema-drift.md"
	// Violations is the governance violations file in .context/state/.
	Violations = "violations.json"
	// TempSuffixPattern is the os.CreateTemp pattern suffix appended
	// to the dot-prefixed base name when staging an atomic write
	// (e.g. ".opencode.json.tmp.*"). The trailing "*" is replaced by
	// CreateTemp with a random unique token.
	TempSuffixPattern = ".tmp.*"
)
