//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package placeholders reads the embedded placeholder
// locale files used by `internal/validate.RejectPlaceholder`
// to filter out lazy body-flag values (TBD, n/a, see chat,
// etc.) on `ctx decision add` / `ctx learning add`.
//
// # Public Surface
//
//   - **[Load](locale)**: returns the folded placeholder
//     set for a locale. Memoized per locale; safe for
//     concurrent use.
//   - **[Reset]**: clears the in-process cache. Test-only.
//
// # Locale Files
//
// Locale YAML lives under
// `internal/assets/i18n/placeholders/<locale>.yaml` and
// is embedded via `assets.FS`. Only `en.yaml` ships in
// v1; the directory layout is stable so additional
// locales (`tr.yaml`, etc.) drop in without code changes.
//
// # Folding Contract
//
// Keys in the returned set are already case-folded via
// [internal/i18n.Fold] at load time. Callers compare
// against `i18n.Fold(strings.TrimSpace(input))`. The
// validator hot path therefore folds once per call and
// does a single O(1) map lookup.
package placeholders
