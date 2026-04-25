//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package io

import "os"

// AppendBytes opens path in append mode, writes data, and closes.
// Returns the first non-nil error encountered among open, write, and
// close. Callers decide whether to propagate, log, or absorb.
//
// Previously this helper logged errors to stderr and returned void
// (best-effort), which conflated "the write succeeded" with "the
// write failed but you'll only know if you scroll stderr". Audit
// trails that depend on the append landing (event.Append, stat
// rollups) need the error to propagate so callers can honour a
// log-first ordering: if the record can't be written, downstream
// side effects should not pretend the event happened.
//
// Parameters:
//   - path: file path to append to (created if missing)
//   - data: bytes to append
//   - perm: file permission bits for creation
//
// Returns:
//   - error: non-nil on open, write, or close failure. When write
//     succeeds but close fails, the close error is returned so
//     disk-flush / fsync problems surface.
func AppendBytes(path string, data []byte, perm os.FileMode) error {
	f, openErr := SafeAppendFile(path, perm)
	if openErr != nil {
		return openErr
	}
	_, writeErr := f.Write(data)
	closeErr := f.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}
