//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package jsonpayload decodes the --json-file argument for the add
// command and overlays its fields onto the noun's flags and content.
//
// # Why
//
// The canonical permissions.deny set matches on the literal Bash
// command string, including the content of --rationale/--context/
// --consequence values. A JSON payload file keeps those values off the
// command line, so a legitimate value that happens to contain a denied
// substring no longer trips the deny rule.
//
// # Shape
//
// [Load] strictly decodes a single JSON object into a [Payload].
// Unknown keys are an error so typos surface instead of silently
// dropping a field. Every key is optional; each add noun consumes only
// the fields relevant to it (extra keys decode without error and are
// ignored by that noun's formatter). Provenance may be supplied on the
// command line or folded into the "provenance" envelope.
//
// # Overlay
//
// The payload reaches the pipeline through two explicit touchpoints:
//
//  1. [OverlayFlags] runs in each noun's PreRunE. It writes the
//     payload's non-empty typed fields (context, rationale, …,
//     provenance) onto the cobra flags via flags.Set, so the
//     decision/learning placeholder gate validates the effective
//     values and run.Run sees them through the bound flag variables.
//     JSON values supersede individually-supplied flags.
//  2. [Payload.Content] supplies the entry content. The extract
//     subpackage calls it first, ahead of --file/args/stdin, when
//     AddConfig.JSONFile is set.
//
// Both touchpoints load the (small) file independently; OverlayFlags
// handles the typed flags, extract handles the positional content.
package jsonpayload
