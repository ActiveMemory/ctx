//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package version parses semantic versions and checks
// encryption key age for the check-version hook. It
// provides helpers that the hook layer calls to build
// nudge boxes when action is needed.
//
// # Semantic Version Parsing
//
// [ParseMajorMinor] extracts the major and minor version
// numbers from a semver string like "1.2.3". Returns
// ok=false for unparseable versions. The hook layer uses
// this to compare the running version against the latest
// available release.
//
// # Key Rotation Check
//
// [CheckKeyAge] checks whether the encryption key at
// rc.KeyPath is older than the configured rotation
// threshold (rc.KeyRotationDays). When the key exceeds
// the threshold:
//
//  1. A nudge message is loaded from the hook template
//     system, with the key age in days as a variable
//  2. The message is wrapped in a NudgeBox with a relay
//     prefix for the agent session
//  3. A notification is emitted via the notify package
//     and relayed to the session
//
// Returns an empty string when the key is fresh or does
// not exist.
package version
