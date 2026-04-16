//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package schema defines the typed error constructors
// and sentinel errors for schema validation. The
// primary export is the [ErrDrift] sentinel, which
// signals that JSONL schema drift was detected.
//
// # Domain
//
// A single sentinel and its constructor cover the
// entire surface:
//
//   - [ErrDrift]: a package-level sentinel error
//     variable. Callers can match it with
//     errors.Is(err, schema.ErrDrift).
//   - [Drift]: convenience constructor that
//     returns the ErrDrift sentinel.
//
// The schema check command and the journal import
// pipeline both return ErrDrift when JSONL fields
// do not match the expected schema. Drift warnings
// are informational; they trigger a non-zero
// exit code but never block operations.
//
// # Wrapping Strategy
//
// ErrDrift is a plain errors.New sentinel with
// no cause wrapping. Its message text comes from
// [internal/config/schema.ErrMsgDrift].
//
// # Concurrency
//
// The sentinel is a package-level variable
// initialized at import time. Safe for concurrent
// use.
package schema
