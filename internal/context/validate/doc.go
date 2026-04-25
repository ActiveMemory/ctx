//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package validate checks whether the .context/
// directory exists and contains all required files.
//
// # Initialization Check
//
// Initialized reports whether the context directory
// contains every file listed in ctx.FilesRequired.
// It stats each file and returns false if any are
// missing. This is used by the init command to decide
// whether to scaffold new files or skip existing ones.
//
//	if validate.Initialized(contextDir) {
//	    // already initialized
//	}
//
// # Existence Check
//
// Exists checks whether a context directory exists on
// disk and is a directory (not a file). If an empty
// string is passed, it falls back to the configured
// context directory from the rc package.
//
//	exists, err := validate.Exists("")
//	if err != nil {
//	    return err
//	}
//	if exists {
//	    // default context dir exists
//	}
//	exists, err = validate.Exists("/custom/path")
//	if err != nil {
//	    return err
//	}
//	if exists {
//	    // custom path exists
//	}
//
// # Usage Pattern
//
// Most commands call Exists first to verify a context
// directory is present, then Initialized to verify it
// has been scaffolded. Commands that modify context
// files guard writes behind both checks.
package validate
