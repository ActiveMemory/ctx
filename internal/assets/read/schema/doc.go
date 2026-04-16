//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package schema provides access to the embedded JSON
// Schema for .ctxrc validation.
//
// # Schema Access
//
// Schema returns the raw ctxrc.schema.json bytes from
// the embedded filesystem. The schema documents all
// configuration fields, their types, defaults, and
// constraints for the .ctxrc configuration file.
//
//	data, err := schema.Schema()
//
// # Validation
//
// The returned JSON Schema can be used by editors and
// tools to validate .ctxrc files. It covers all
// fields in the CtxRC struct and is kept in sync via
// TestSchemaCoversCtxRC.
//
// # Sync Guarantee
//
// The schema is tested against the CtxRC struct to
// ensure every field is documented. If a new field is
// added to the struct without updating the schema,
// the test fails, preventing drift.
package schema
