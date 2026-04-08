//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings and DescKeys for the connect command group.
const (
	// UseConnect is the cobra Use string for connect.
	UseConnect = "connect"
	// UseConnectRegister is the Use string for register.
	UseConnectRegister = "register <hub-address>"
	// UseConnectSubscribe is the Use string for subscribe.
	UseConnectSubscribe = "subscribe <types...>"
	// UseConnectSync is the Use string for sync.
	UseConnectSync = "sync"
	// UseConnectPublish is the Use string for publish.
	UseConnectPublish = "publish"
	// UseConnectListen is the Use string for listen.
	UseConnectListen = "listen"
	// UseConnectStatus is the Use string for status.
	UseConnectStatus = "status"

	// DescKeyConnect is the desc key for the connect command.
	DescKeyConnect = "connect"
	// DescKeyConnectRegister is the desc key for register.
	DescKeyConnectRegister = "connect.register"
	// DescKeyConnectSubscribe is the desc key for subscribe.
	DescKeyConnectSubscribe = "connect.subscribe"
	// DescKeyConnectSync is the desc key for sync.
	DescKeyConnectSync = "connect.sync"
	// DescKeyConnectPublish is the desc key for publish.
	DescKeyConnectPublish = "connect.publish"
	// DescKeyConnectListen is the desc key for listen.
	DescKeyConnectListen = "connect.listen"
	// DescKeyConnectStatus is the desc key for status.
	DescKeyConnectStatus = "connect.status"
)
