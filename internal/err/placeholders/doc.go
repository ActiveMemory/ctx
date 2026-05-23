//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package placeholders provides typed errors for failures
// in `internal/assets/read/placeholders` — the loader for
// the embedded placeholder locale YAML.
//
// # Public Surface
//
//   - **[ReadLocale](locale, cause)**: file-read failure
//     wrapping the locale identifier and the underlying
//     OS error.
//   - **[ParseLocale](locale, cause)**: YAML-parse failure
//     wrapping the locale identifier and the underlying
//     yaml.v3 error.
//
// Both error functions are build-time invariant violations
// (the YAML is embedded in the binary). They exist for
// telemetry / debug surfacing; the validator that consumes
// the loader fails closed when either is returned.
package placeholders
