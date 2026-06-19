//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

// ensure initializes the registry's maps for zero-value use.
//
// Parameters:
//   - registry: registry receiver to initialize
func (registry *Registry) ensure() {
	if registry.factories == nil {
		registry.factories = make(map[string]Factory)
	}
	if registry.configs == nil {
		registry.configs = make(map[string]Config)
	}
}
