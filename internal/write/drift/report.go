//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// ReportHeader prints the drift report heading, separator, and trailing
// blank line. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
func ReportHeader(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyDriftReportHeading))
	cmd.Println(desc.Text(text.DescKeyDriftReportSeparator))
	cmd.Println()
}

// ViolationsHeading prints the violations section heading with count
// and a trailing blank line. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: number of violations
func ViolationsHeading(cmd *cobra.Command, count int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyDriftViolationsHeading), count))
	cmd.Println()
}

// ViolationLine prints a single violation entry. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - line: pre-formatted violation line
func ViolationLine(cmd *cobra.Command, line string) {
	if cmd == nil {
		return
	}
	cmd.Println(line)
}

// WarningsHeading prints the warnings section heading with count
// and a trailing blank line. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: number of warnings
func WarningsHeading(cmd *cobra.Command, count int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyDriftWarningsHeading), count))
	cmd.Println()
}

// PathRefsBlock prints the path references category label followed by
// each item line and a trailing blank line. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - items: pre-formatted path reference lines
func PathRefsBlock(cmd *cobra.Command, items []string) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyDriftPathRefsLabel))
	for _, item := range items {
		cmd.Println(item)
	}
	cmd.Println()
}

// StalenessBlock prints the staleness category label followed by
// each item line and a trailing blank line. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - items: pre-formatted staleness lines
func StalenessBlock(cmd *cobra.Command, items []string) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyDriftStalenessLabel))
	for _, item := range items {
		cmd.Println(item)
	}
	cmd.Println()
}

// OtherBlock prints the other warnings category label followed by
// each item line and a trailing blank line. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - items: pre-formatted other warning lines
func OtherBlock(cmd *cobra.Command, items []string) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyDriftOtherLabel))
	for _, item := range items {
		cmd.Println(item)
	}
	cmd.Println()
}

// PassedHeading prints the passed checks heading with count.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: number of passed checks
func PassedHeading(cmd *cobra.Command, count int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyDriftPassedHeading), count))
}

// PassedLine prints a single passed check entry. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: formatted check name
func PassedLine(cmd *cobra.Command, name string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyDriftPassedLine), name))
}

// StatusViolation prints the violation status verdict with leading
// blank line. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
func StatusViolation(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyDriftStatusViolation))
}

// StatusWarning prints the warning status verdict with leading
// blank line. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
func StatusWarning(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyDriftStatusWarning))
}

// StatusOK prints the OK status verdict with leading blank line.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
func StatusOK(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyDriftStatusOK))
}

// BlankLine prints a blank line. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
func BlankLine(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println()
}
