//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package catalog maps context files to MCP resource
// URIs and builds the resource list returned by the
// resources/list method.
//
// # Initialization
//
// Init must be called once during server bootstrap
// before FileForURI is used. It populates an internal
// lookup map from the static resource table.
//
//	catalog.Init()
//
// # URI Construction
//
// URI builds a full resource URI from a name suffix
// by prepending the configured URI prefix:
//
//	catalog.URI("tasks")
//	// => "ctx://context/tasks"
//
// AgentURI returns the URI for the assembled agent
// packet, which combines all context files into a
// single response.
//
// # Lookup
//
// FileForURI returns the context file name for a
// given resource URI, or empty string if the URI
// does not correspond to a known file resource.
//
//	name := catalog.FileForURI(uri)
//
// # Resource List
//
// ToList constructs the immutable resource list that
// the server returns for resources/list requests. It
// includes all individual file resources plus the
// agent packet resource.
//
// # Types
//
// The mapping type pairs a context file name with its
// MCP resource name and human-readable description.
// The static table variable holds all known resources.
package catalog
