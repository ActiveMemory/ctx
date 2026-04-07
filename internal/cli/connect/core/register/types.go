//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package register

// connConfig is the persisted connection configuration.
//
// Fields:
//   - HubAddr: hub gRPC address
//   - Token: client bearer token
type connConfig struct {
	HubAddr string `json:"hub_addr"`
	Token   string `json:"token"`
}
