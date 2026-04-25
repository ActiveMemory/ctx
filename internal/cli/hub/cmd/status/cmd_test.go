//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/cli"
)

// TestHubStatus_AnnotationSkipInit guards the hub-bypass contract.
// Spec: specs/single-source-context-anchor.md.
func TestHubStatus_AnnotationSkipInit(t *testing.T) {
	c := Cmd()
	if got, ok := c.Annotations[cli.AnnotationSkipInit]; !ok {
		t.Errorf("hub status: missing AnnotationSkipInit annotation")
	} else if got != cli.AnnotationTrue {
		t.Errorf("hub status: AnnotationSkipInit = %q, want %q", got, cli.AnnotationTrue)
	}
}
