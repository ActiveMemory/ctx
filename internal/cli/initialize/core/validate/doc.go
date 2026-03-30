//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package validate runs pre-flight checks during initialization.
//
// [CheckCtxInPath] verifies the ctx binary is accessible via PATH,
// warning the user if the installed binary will not be found by
// hooks and scripts.
package validate
