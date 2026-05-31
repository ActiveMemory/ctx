//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tpl

// ZensicalProject is the static [project] section of the generated
// zensical.toml. It has no interpolation; it is loaded verbatim from
// an embedded file at init and written through by callers as-is.
var ZensicalProject string

// ZensicalTheme is the static theme and extras section of the
// generated zensical.toml, loaded verbatim from an embedded file at
// init.
var ZensicalTheme string
