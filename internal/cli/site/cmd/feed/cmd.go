//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package feed provides the "ctx site feed" subcommand.
package feed

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/rss"
	"github.com/ActiveMemory/ctx/internal/flagbind"
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

	c := &cobra.Command{
		Use:     cmd.UseSiteFeed,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeySiteFeed),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, rss.DefaultFeedInputDir, out, baseURL)
		},
	}

	flagbind.StringFlagPDefault(c, &out,
		cFlag.Out, cFlag.ShortOutput,
		rss.DefaultFeedOutPath,
		flag.DescKeySiteFeedOut,
	)
	flagbind.StringFlagDefault(c, &baseURL,
		cFlag.BaseURL, rss.DefaultFeedBaseURL,
		flag.DescKeySiteFeedBaseUrl,
	)

	return c
}
