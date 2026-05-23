//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/entity"
)

const (
	// ErrBackendNotFound signals a Resolve call naming a
	// backend that has no Config or no Factory registered.
	// Constructor [NotFound] wraps it with the requested
	// name.
	ErrBackendNotFound = entity.Sentinel(
		text.DescKeyErrBackendBackendNotFoundMsg,
	)
	// ErrNoBackends signals a Resolve or Default call
	// against a Registry with no backends configured at
	// all. Recovery: `ctx setup --backend <name>`.
	ErrNoBackends = entity.Sentinel(
		text.DescKeyErrBackendNoBackends,
	)
	// ErrAmbiguousDefault signals a Default call against a
	// Registry where more than one backend is configured
	// and no default has been set via SetDefault.
	// Recovery: pass `--backend <name>` or set
	// `[backends].default` in `.ctxrc`.
	ErrAmbiguousDefault = entity.Sentinel(
		text.DescKeyErrBackendAmbiguousDefault,
	)
	// ErrDuplicateRegistration signals a Register call for
	// a type name already bound to a Factory. Constructor
	// [DuplicateRegistration] wraps it with the offending
	// name.
	ErrDuplicateRegistration = entity.Sentinel(
		text.DescKeyErrBackendDuplicateRegistrationMsg,
	)
)

// NotFound wraps [ErrBackendNotFound] with the requested
// backend name so error output names the missing key.
//
// Parameters:
//   - name: the backend type label the caller asked for.
//
// Returns:
//   - error: wraps [ErrBackendNotFound] for errors.Is.
func NotFound(name string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackendBackendNotFound),
		ErrBackendNotFound, name,
	)
}

// DuplicateRegistration wraps [ErrDuplicateRegistration]
// with the offending type name so error output identifies
// the conflict.
//
// Parameters:
//   - name: the type label that was already bound.
//
// Returns:
//   - error: wraps [ErrDuplicateRegistration] for
//     errors.Is.
func DuplicateRegistration(name string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrBackendDuplicateRegistration),
		ErrDuplicateRegistration, name,
	)
}
