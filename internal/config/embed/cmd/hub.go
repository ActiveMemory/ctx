//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings and DescKeys for the hub command group.
const (
	// UseHub is the cobra Use string for hub.
	UseHub = "hub"
	// UseHubStatus is the Use string for hub status.
	UseHubStatus = "status"
	// UseHubPeer is the Use string for hub peer.
	UseHubPeer = "peer <add|remove> <address>"
	// UseHubStepdown is the Use string for hub stepdown.
	UseHubStepdown = "stepdown"

	// DescKeyHub is the desc key for the hub command.
	DescKeyHub = "hub"
	// DescKeyHubStatus is the desc key for hub status.
	DescKeyHubStatus = "hub.status"
	// DescKeyHubPeer is the desc key for hub peer.
	DescKeyHubPeer = "hub.peer"
	// DescKeyHubStepdown is the desc key for hub stepdown.
	DescKeyHubStepdown = "hub.stepdown"
)
