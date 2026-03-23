//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package bootstrap initializes the ctx CLI application.
//
// It provides functions to create the root command and register all
// subcommands. The typical usage pattern is:
//
//	cmd := bootstrap.Initialize(bootstrap.RootCmd())
//	if err := cmd.Execute(); err != nil {
//	    // handle error
package bootstrap
