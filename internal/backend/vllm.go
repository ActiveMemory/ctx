//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import "context"

// Ping implements [Backend] with vLLM-specific cold-start
// retry. Delegates to the embedded openAICompat.Ping via
// [coldStartRetry] which retries while the failure is a
// dial-refused error (server has not yet bound the
// listener) and the cold-start window has not elapsed.
//
// Parameters:
//   - ctx: caller-provided context for cancellation.
//
// Returns:
//   - error: nil on success, typed sentinel on failure.
func (b *vllm) Ping(ctx context.Context) error {
	return coldStartRetry(
		ctx,
		b.openAICompat.Ping,
		b.coldStartWindow,
		b.coldStartInterval,
	)
}
