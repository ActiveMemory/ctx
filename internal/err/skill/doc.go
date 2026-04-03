//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package skill provides error constructors for skill operations.
//
// Error constructors return structured errors with context for
// user-facing messages routed through internal/assets text lookups.
// Exports: [CreateDest], [Install], [InvalidManifest],
// [InvalidYAML], [List], [Load], [MissingClosingDelimiter],
// [MissingName], [MissingOpeningDelimiter], [NotFound],
// [NotValidDir], [NotValidSource], [Read], [ReadDir],
// [Remove], [SkillLoad].
package skill
