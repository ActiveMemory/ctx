//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for AI-backend registry errors. The matching
// YAML entries live in commands/text/errors.yaml under
// the `err.backend.*` namespace; sentinels and
// constructors in internal/err/backend/ resolve them via
// desc.Text at error construction time.
const (
	// DescKeyErrBackendBackendNotFoundMsg is the sentinel
	// message lookup key for ErrBackendNotFound.
	DescKeyErrBackendBackendNotFoundMsg = "err.backend.backend-not-found-msg"
	// DescKeyErrBackendBackendNotFound is the wrapper
	// format key for the BackendNotFound constructor;
	// interpolates the sentinel and the requested name.
	DescKeyErrBackendBackendNotFound = "err.backend.backend-not-found"
	// DescKeyErrBackendNoBackends is the sentinel message
	// key for ErrNoBackends (no wrapper; the sentinel
	// itself carries the user-facing recovery hint).
	DescKeyErrBackendNoBackends = "err.backend.no-backends"
	// DescKeyErrBackendAmbiguousDefault is the sentinel
	// message key for ErrAmbiguousDefault.
	DescKeyErrBackendAmbiguousDefault = "err.backend.ambiguous-default"
	// DescKeyErrBackendDuplicateRegistrationMsg is the
	// sentinel message lookup key for
	// ErrDuplicateRegistration.
	DescKeyErrBackendDuplicateRegistrationMsg = "err.backend.duplicate-registration-msg"
	// DescKeyErrBackendDuplicateRegistration is the
	// wrapper format key for the DuplicateRegistration
	// constructor.
	DescKeyErrBackendDuplicateRegistration = "err.backend.duplicate-registration"
)
