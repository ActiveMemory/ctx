//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sanitize converts internal drift check
// identifiers into human-readable display names.
//
// Drift checks are identified by machine-friendly constants
// such as "path_references" or "staleness_check" defined in
// [internal/config/drift]. These identifiers are stable for
// configuration and code references but are not suitable for
// end-user output.
//
// [FormatCheckName] maps each known [cfgDrift.CheckName] to
// a descriptive label loaded from the embedded text asset
// system via [desc.Text]. Unknown identifiers pass through
// unchanged, so new checks degrade gracefully before their
// labels are added to the asset YAML.
//
// # Design Choice
//
// Labels are resolved at call time from the asset cache
// rather than hardcoded in a switch, so the same text
// definition is shared between the drift CLI output and
// any other consumer that needs check names (e.g., doctor
// checks, hook messages).
package sanitize
