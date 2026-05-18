//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for schema validation error sentinels. The matching
// YAML entry lives in commands/text/errors.yaml; constructors in
// internal/err/schema/ resolve it via desc.Text at error
// construction time.
const (
	// DescKeyErrSchemaDrift is the text key for the
	// schema-drift sentinel.
	DescKeyErrSchemaDrift = "err.schema.drift"
)
