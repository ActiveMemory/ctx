//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package opencode

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	cfgFS "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgToken "github.com/ActiveMemory/ctx/internal/config/token"
	errFS "github.com/ActiveMemory/ctx/internal/err/fs"
	setupErr "github.com/ActiveMemory/ctx/internal/err/setup"
	ctxio "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// EnsureProviderConfig merges backend baseURL wiring into the user's
// OpenCode global config.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - backendName: configured ctx backend name
//   - endpoint: resolved backend endpoint URL
//
// Returns:
//   - error: read, merge, marshal, or write failure
func EnsureProviderConfig(
	cmd *cobra.Command,
	backendName string,
	endpoint string,
) error {
	target, pathErr := globalConfigPath()
	if pathErr != nil {
		return pathErr
	}
	if _, validateErr := validateManagedTarget(target); validateErr != nil {
		return validateErr
	}
	providerID, npmPackage, ok := providerDetails(backendName)
	if !ok {
		return nil
	}
	config := map[string]any{}
	data, readErr := ctxio.SafeReadUserFile(target)
	if readErr == nil {
		if len(data) > 0 {
			if unmarshalErr := json.Unmarshal(data, &config); unmarshalErr != nil {
				return setupErr.MarshalConfig(unmarshalErr)
			}
		}
	} else if !os.IsNotExist(readErr) {
		return errFS.FileRead(target, readErr)
	}
	mergeProviderConfig(config, providerID, npmPackage, endpoint)
	out, marshalErr := json.MarshalIndent(config, "", cfgToken.Indent2)
	if marshalErr != nil {
		return setupErr.MarshalConfig(marshalErr)
	}
	out = append(out, cfgToken.NewlineLF...)
	if writeFileErr := ctxio.SafeWriteFileAtomic(
		target,
		out,
		cfgFS.PermFile,
	); writeFileErr != nil {
		return errFS.FileWrite(target, writeFileErr)
	}
	writeSetup.InfoOpenCodeCreated(cmd, target)
	return nil
}
