//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

// Error message constants for rc sentinel errors. These are used
// only for errors.Is matching; user-facing wrapping goes through
// err/context constructors that format tailored messages.
const (
	// ErrMsgDirNotDeclared is the sentinel message for the
	// "context directory has not been declared" error.
	ErrMsgDirNotDeclared = "context directory not declared"
	// ErrMsgRelativeNotAllowed is the sentinel message for the
	// "CTX_DIR must be absolute" rejection.
	ErrMsgRelativeNotAllowed = "context directory must be absolute"
	// ErrMsgNonCanonicalBasename is the sentinel message for the
	// "CTX_DIR basename must be .context" rejection.
	ErrMsgNonCanonicalBasename = "context directory has non-canonical basename"
	// ErrMsgContextDirNotFound is the sentinel message for the
	// "declared CTX_DIR does not exist" rejection.
	ErrMsgContextDirNotFound = "context directory not found: "
	// ErrMsgContextDirNotADirectory is the sentinel message for the
	// "CTX_DIR is a file, not a directory" rejection.
	ErrMsgContextDirNotADirectory = "context directory is not a directory"
	// ErrMsgContextDirStat is the sentinel message for stat failures
	// other than not-exist (permission denied, I/O error).
	ErrMsgContextDirStat = "context directory stat failed"
)

// Format strings for sentinel-wrapping in err/context constructors.
// Centralized here so the magic-string audit (which exempts
// internal/config) does not flag them at the call site.
const (
	// FmtWrapColon wraps a sentinel and a tailored message:
	//   fmt.Errorf(FmtWrapColon, ErrFoo, "tailored detail")
	//   ↦ "<ErrFoo.Error()>: tailored detail".
	FmtWrapColon = "%w: %s"
	// FmtWrapBare appends the tailored detail directly to the
	// sentinel without a separator. Used when the sentinel message
	// already ends with whatever separator the caller wants
	// (e.g., a trailing space-colon for "context directory not found: ").
	FmtWrapBare = "%w%s"
)
