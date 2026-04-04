//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings for steering subcommands.
const (
	// UseSteering is the cobra Use string for the steering command.
	UseSteering = "steering"
	// UseSteeringAdd is the cobra Use string for the steering add command.
	UseSteeringAdd = "add <name>"
	// UseSteeringList is the cobra Use string for the steering list command.
	UseSteeringList = "list"
	// UseSteeringPreview is the cobra Use string for the steering preview command.
	UseSteeringPreview = "preview <prompt>"
	// UseSteeringInit is the cobra Use string for the steering init command.
	UseSteeringInit = "init"
	// UseSteeringSync is the cobra Use string for the steering sync command.
	UseSteeringSync = "sync"
)

// DescKeys for steering subcommands.
const (
	// DescKeySteering is the description key for the steering command.
	DescKeySteering = "steering"
	// DescKeySteeringAdd is the description key for the steering add command.
	DescKeySteeringAdd = "steering.add"
	// DescKeySteeringList is the description key for the steering list command.
	DescKeySteeringList = "steering.list"
	// DescKeySteeringPreview is the description key for the steering preview
	// command.
	DescKeySteeringPreview = "steering.preview"
	// DescKeySteeringInit is the description key for the steering init command.
	DescKeySteeringInit = "steering.init"
	// DescKeySteeringSync is the description key for the steering sync command.
	DescKeySteeringSync = "steering.sync"
)
