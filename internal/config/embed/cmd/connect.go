//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings and DescKeys for the connection command group.
const (
	// UseConnection is the cobra Use string for connection.
	UseConnection = "connection"
	// UseConnectionRegister is the Use string for register.
	UseConnectionRegister = "register <hub-address>"
	// UseConnectionSubscribe is the Use string for subscribe.
	UseConnectionSubscribe = "subscribe <types...>"
	// UseConnectionSync is the Use string for sync.
	UseConnectionSync = "sync"
	// UseConnectionPublish is the Use string for publish.
	UseConnectionPublish = "publish"
	// UseConnectionListen is the Use string for listen.
	UseConnectionListen = "listen"
	// UseConnectionStatus is the Use string for status.
	UseConnectionStatus = "status"

	// DescKeyConnection is the desc key for the connection command.
	DescKeyConnection = "connection"
	// DescKeyConnectionRegister is the desc key for register.
	DescKeyConnectionRegister = "connection.register"
	// DescKeyConnectionSubscribe is the desc key for subscribe.
	DescKeyConnectionSubscribe = "connection.subscribe"
	// DescKeyConnectionSync is the desc key for sync.
	DescKeyConnectionSync = "connection.sync"
	// DescKeyConnectionPublish is the desc key for publish.
	DescKeyConnectionPublish = "connection.publish"
	// DescKeyConnectionListen is the desc key for listen.
	DescKeyConnectionListen = "connection.listen"
	// DescKeyConnectionStatus is the desc key for status.
	DescKeyConnectionStatus = "connection.status"
)
