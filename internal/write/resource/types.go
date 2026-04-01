//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resource

type jsonAlert struct {
	Severity string `json:"severity"`
	Resource string `json:"resource"`
	Message  string `json:"message"`
}

type jsonOutput struct {
	Memory struct {
		TotalBytes uint64 `json:"total_bytes"`
		UsedBytes  uint64 `json:"used_bytes"`
		Percent    int    `json:"percent"`
		Supported  bool   `json:"supported"`
	} `json:"memory"`
	Swap struct {
		TotalBytes uint64 `json:"total_bytes"`
		UsedBytes  uint64 `json:"used_bytes"`
		Percent    int    `json:"percent"`
		Supported  bool   `json:"supported"`
	} `json:"swap"`
	Disk struct {
		TotalBytes uint64 `json:"total_bytes"`
		UsedBytes  uint64 `json:"used_bytes"`
		Percent    int    `json:"percent"`
		Path       string `json:"path"`
		Supported  bool   `json:"supported"`
	} `json:"disk"`
	Load struct {
		Load1     float64 `json:"load1"`
		Load5     float64 `json:"load5"`
		Load15    float64 `json:"load15"`
		NumCPU    int     `json:"num_cpu"`
		Ratio     float64 `json:"ratio"`
		Supported bool    `json:"supported"`
	} `json:"load"`
	Alerts      []jsonAlert `json:"alerts"`
	MaxSeverity string      `json:"max_severity"`
}
