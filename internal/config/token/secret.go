//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package token

// SecretPatterns are filename substrings that indicate potential secret files.
var SecretPatterns = []string{
	".env",
	"credentials",
	"secret",
	"api_key",
	"apikey",
	"password",
}
