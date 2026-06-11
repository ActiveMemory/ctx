//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package accept wires the "ctx dream accept <id>" subcommand.
//
// It builds the cobra command and delegates to the dispose core
// logic, which loads the proposal by id and applies its recommended
// action through the engine. Mechanical actions complete here;
// generative ones are recorded as intent for /ctx-serendipity.
package accept
