//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cli

// AnnotationSkipInit is the cobra.Command annotation key that exempts
// a command from the PersistentPreRunE initialization guard.
const AnnotationSkipInit = "skipInitCheck"

// AnnotationTrue is the canonical value for boolean cobra annotations.
const AnnotationTrue = "true"
