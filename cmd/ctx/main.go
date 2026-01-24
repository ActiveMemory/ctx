//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// package main is the main entry point of the app.
package main

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/bootstrap"
)

func main() {
	cmd := bootstrap.Initialize(bootstrap.RootCmd())

	if err := cmd.Execute(); err != nil {
		cmd.PrintErrln("Error:", err)
		os.Exit(1)
	}
}
