//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package catalog lists available context template
// files from the embedded assets.
//
// # Listing Templates
//
// List returns the names of all .context/ template
// files (TASKS.md, DECISIONS.md, CONVENTIONS.md, etc.)
// available for deployment by ctx init. The list is
// derived from the embedded filesystem at compile
// time by reading the context/ asset directory.
//
//	names, err := catalog.List()
//	for _, name := range names {
//	    fmt.Println(name) // "TASKS.md", etc.
//	}
//
// # Usage
//
// The init command uses this list to determine which
// template files to deploy into a new .context/
// directory. Each name corresponds to a file that
// can be read via the template package.
package catalog
