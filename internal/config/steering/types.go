//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

// InclusionMode determines when a steering file is
// injected into an AI prompt.
type InclusionMode string

// Inclusion mode constants for steering file injection.
const (
	// InclusionAlways includes the file in every packet.
	InclusionAlways InclusionMode = "always"
	// InclusionAuto includes when prompt matches.
	InclusionAuto InclusionMode = "auto"
	// InclusionManual includes only when named.
	InclusionManual InclusionMode = "manual"
)
