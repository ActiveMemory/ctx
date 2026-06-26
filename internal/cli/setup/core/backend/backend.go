//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"io"
	"os"

	cfgFS "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
	setupErr "github.com/ActiveMemory/ctx/internal/err/setup"
	ctxio "github.com/ActiveMemory/ctx/internal/io"
)

// Run writes or prints backend setup configuration.
//
// Parameters:
//   - out: destination for dry-run output and warnings
//   - options: backend setup options
//
// Returns:
//   - error: read, marshal, or write failure
func Run(out io.Writer, options Options) error {
	if validateErr := validate(options); validateErr != nil {
		return validateErr
	}
	resolved := defaults(options)
	if resolved.APIKeyEnv != "" && os.Getenv(resolved.APIKeyEnv) != "" {
		if _, warnErr := io.WriteString(
			out,
			cfgSetup.BackendEnvWarn+resolved.APIKeyEnv+
				cfgSetup.BackendEnvWarnEnd,
		); warnErr != nil {
			return setupErr.WriteFile(cfgSetup.FileCtxRC, warnErr)
		}
	}
	downstream := downstreamEnv(resolved.Name, resolved.Endpoint)
	if downstream != "" {
		if _, warnErr := io.WriteString(out, downstream); warnErr != nil {
			return setupErr.WriteFile(cfgSetup.FileCtxRC, warnErr)
		}
	}
	content, contentErr := content(resolved)
	if contentErr != nil {
		return contentErr
	}
	if !resolved.Write {
		if _, writeErr := io.WriteString(
			out,
			cfgSetup.BackendDryRunPrefix,
		); writeErr != nil {
			return setupErr.WriteFile(cfgSetup.FileCtxRC, writeErr)
		}
		if _, writeErr := out.Write(content); writeErr != nil {
			return setupErr.WriteFile(cfgSetup.FileCtxRC, writeErr)
		}
		return nil
	}
	if writeErr := ctxio.SafeWriteFile(
		cfgSetup.FileCtxRC,
		content,
		cfgFS.PermFile,
	); writeErr != nil {
		return setupErr.WriteFile(cfgSetup.FileCtxRC, writeErr)
	}
	if _, doneErr := io.WriteString(
		out,
		cfgSetup.BackendWriteDone,
	); doneErr != nil {
		return setupErr.WriteFile(cfgSetup.FileCtxRC, doneErr)
	}
	return nil
}
