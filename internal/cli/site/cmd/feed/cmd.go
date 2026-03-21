//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package feed provides the "ctx site feed" subcommand.
package feed

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/ActiveMemory/ctx/internal/config/rss"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx site feed" subcommand.
//
// Returns:
//   - *cobra.Command: Configured feed generation subcommand
func Cmd() *cobra.Command {
	var (
		out     string
		baseURL string
	)

	short, long := desc.Command(cmd.DescKeySiteFeed)

	cmd := &cobra.Command{
		Use:   "feed",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runFeed(cmd, rss.DefaultFeedInputDir, out, baseURL)
		},
	}

	cmd.Flags().StringVarP(
		&out, "out", "o", rss.DefaultFeedOutPath,
		desc.Flag(flag.DescKeySiteFeedOut),
	)
	cmd.Flags().StringVar(
		&baseURL, "base-url", rss.DefaultFeedBaseURL,
		desc.Flag(flag.DescKeySiteFeedBaseUrl),
	)

	return cmd
}
