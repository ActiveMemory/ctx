//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package admin

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/config/env"
	errHub "github.com/ActiveMemory/ctx/internal/err/hub"
)

// Token resolves the hub admin credential.
//
// The --token flag takes precedence, then the
// CTX_HUB_ADMIN_TOKEN environment variable.
//
// Parameters:
//   - flagValue: value of the command's --token flag
//
// Returns:
//   - string: the resolved admin token
//   - error: non-nil when neither source supplied one
func Token(flagValue string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}

	if fromEnv := os.Getenv(env.HubAdmin); fromEnv != "" {
		return fromEnv, nil
	}

	return "", errHub.AdminTokenRequired()
}
