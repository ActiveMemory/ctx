//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for ctx-init kb-scaffolding error wrappers. The
// matching YAML entries live in commands/text/errors.yaml;
// constructors in internal/err/initialize/kb/ resolve them via
// desc.Text at error construction time.
const (
	// DescKeyErrInitKbMkdir wraps `os.MkdirAll` for a scaffolded
	// directory.
	DescKeyErrInitKbMkdir = "err.initkb.mkdir"
	// DescKeyErrInitKbCopyIngestTemplates wraps an ingest-templates
	// copy failure.
	DescKeyErrInitKbCopyIngestTemplates = "err.initkb.copy-ingest-templates"
	// DescKeyErrInitKbCopySchemas wraps a schemas copy failure.
	DescKeyErrInitKbCopySchemas = "err.initkb.copy-schemas"
	// DescKeyErrInitKbCopyLanding wraps the kb landing-page copy
	// failure.
	DescKeyErrInitKbCopyLanding = "err.initkb.copy-landing"
	// DescKeyErrInitKbReadEmbed wraps an embedded-asset read
	// failure.
	DescKeyErrInitKbReadEmbed = "err.initkb.read-embed"
	// DescKeyErrInitKbMkdirFor wraps a parent-dir create failure
	// for a destination file.
	DescKeyErrInitKbMkdirFor = "err.initkb.mkdir-for"
	// DescKeyErrInitKbWriteFile wraps a destination-file write
	// failure.
	DescKeyErrInitKbWriteFile = "err.initkb.write-file"
	// DescKeyErrInitKbReadDir wraps an `os.ReadDir` failure for a
	// scaffolded directory.
	DescKeyErrInitKbReadDir = "err.initkb.read-dir"
)
