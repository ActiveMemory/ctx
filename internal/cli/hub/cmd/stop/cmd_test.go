//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stop

import (
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/cli"
)

// TestHubStop_AnnotationSkipInit guards the hub-bypass contract.
// Spec: specs/single-source-context-anchor.md.
func TestHubStop_AnnotationSkipInit(t *testing.T) {
	c := Cmd()
	if got, ok := c.Annotations[cli.AnnotationSkipInit]; !ok {
		t.Errorf("hub stop: missing AnnotationSkipInit annotation")
	} else if got != cli.AnnotationTrue {
		t.Errorf("hub stop: AnnotationSkipInit = %q, want %q", got, cli.AnnotationTrue)
	}
}
