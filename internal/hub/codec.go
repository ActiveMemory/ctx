//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import "google.golang.org/grpc/encoding"

// init registers the JSON codec with the gRPC runtime so
// plain Go structs can be used as RPC messages.
func init() { encoding.RegisterCodec(jsonCodec{}) }
