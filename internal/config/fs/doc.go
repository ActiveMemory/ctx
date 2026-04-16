//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package fs defines file and directory permission
// constants used whenever ctx creates files, hooks,
// directories, or secret material on disk.
//
// Unix file permissions are easy to get wrong: a hook
// script missing the execute bit fails silently, and a
// key file with world-readable permissions leaks
// secrets. This package provides named constants so
// callers express intent ("this is a secret file")
// rather than raw octal values.
//
// # File Permissions
//
//   - PermFile (0644): standard regular files;
//     owner read-write, group and others read-only.
//     Used for context markdown files, config files,
//     and generated output.
//   - PermExec (0755): directories and executable
//     files; owner full access, group and others
//     read-execute. Used for .context/ directories
//     and hook scripts.
//   - PermRestrictedDir (0750): internal directories
//     where others should have no access; owner full,
//     group read-execute. Used for state/ and logs/.
//   - PermSecret (0600): secret files; owner
//     read-write only. Used for encryption keys and
//     token files.
//   - PermKeyDir (0700): key storage directories;
//     owner-only full access. Used for the user-level
//     key directory.
//
// # Permission Bit Masks
//
//   - ExecBitMask (0111): bitmask for testing
//     whether any executable bit is set. Used by the
//     drift subsystem to detect hook scripts missing
//     the execute permission.
//
// # Why Centralized
//
// Permission constants are referenced by file writers,
// directory creators, hook installers, key management,
// and drift detection. Defining them in one package
// ensures consistent security posture and makes
// permission auditing straightforward.
package fs
