//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

// Config is the persisted hub connection configuration.
//
// Fields:
//   - HubAddr: hub gRPC address (host:port)
//   - Token: client bearer token for RPCs
//   - Types: subscribed entry types (empty = all)
type Config struct {
	HubAddr string   `json:"hub_addr"`
	Token   string   `json:"token"`
	Types   []string `json:"types,omitempty"`
}
