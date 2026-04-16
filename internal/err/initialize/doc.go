//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package initialize defines the typed error
// constructors for the `ctx init` command and
// bootstrap sequence. These errors fire when
// creating the .context/ directory, deploying
// templates, or verifying preconditions.
//
// # Domain
//
// Errors fall into four categories:
//
//   - **Not initialized**: the project has no
//     .context/ directory. Constructors:
//     [NotInitialized], [ContextNotInitialized].
//   - **Environment**: the home directory cannot
//     be resolved or ctx is not on PATH.
//     Constructors: [HomeDir], [CtxNotInPath].
//   - **Template IO**: an embedded template or
//     project README could not be read, or the
//     Makefile could not be created.
//     Constructors: [ReadTemplate],
//     [ReadProjectReadme], [CreateMakefile].
//   - **Deployment**: listing or reading
//     embedded files during template deployment
//     failed. Constructors: [DeployList],
//     [DeployRead], [DetectReferenceTime].
//
// # Wrapping Strategy
//
// IO constructors wrap their cause with
// fmt.Errorf %w so callers can errors.Is against
// system errors. [NotInitialized] and
// [ContextNotInitialized] return plain errors
// because they signal a missing precondition,
// not an IO failure. All user-facing text is
// resolved through [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package initialize
