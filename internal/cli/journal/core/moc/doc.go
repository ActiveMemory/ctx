//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package moc generates **Maps of Content** — the navigational
// index pages that sit at the top of both the journal site and
// the Obsidian vault and tell a human "here are the high-level
// topics, here are the key files, here are the recent entries
// that matter most".
//
// "MOC" is borrowed from the personal-knowledge-management
// world (Obsidian / Linking Your Thinking) where it names the
// curated dashboard page that aggregates by topic and key
// entity rather than by chronology.
//
// # The Surface
//
//   - **[Home](entries, opts)** — generates the **site
//     homepage MOC**: top topics, key files, recent
//     entries, all in a single page. Output is markdown
//     ready to land at `site/index.md`.
//   - **[ObsidianTopics](entries)** — generates the
//     Obsidian-vault topics index using `[[wikilink]]`
//     syntax. Lives at `vault/MOC.md`.
//   - **[GenerateObsidianTopicPage](topic, entries)** —
//     generates a per-topic page in Obsidian format with
//     wikilinks back to each matching entry. Lives at
//     `vault/topics/<slug>.md`.
//
// # Site MOC vs Obsidian MOC
//
// The two flavors share the *aggregation logic* (topic
// counts, key-file detection, recency ranking) but
// diverge in **link syntax**:
//
//   - The site uses standard `[text](url.md)` markdown
//     links so zensical can resolve them through its
//     navigation graph.
//   - Obsidian uses `[[wikilinks]]` so its native graph
//     view picks them up.
//
// Each helper assembles the link in the right dialect;
// the aggregation results are reused.
//
// # Inputs
//
// All MOC generators take a slice of [entity.Entry] and
// optionally a [TopicIndex] (built by
// [internal/cli/journal/core/section.BuildTopicIndex]).
// The MOC is a *projection* of the entry set, not a
// transformation: original entries are unchanged.
//
// # Concurrency
//
// All functions are pure data transformations over
// the entry slice. Concurrent callers never race.
//
// # Related Packages
//
//   - [internal/cli/journal/cmd/site]      — invokes
//     [Home] when building the zensical site.
//   - [internal/cli/journal/cmd/obsidian]  — invokes
//     [ObsidianTopics] / [GenerateObsidianTopicPage]
//     when exporting the vault.
//   - [internal/cli/journal/core/section]  — produces
//     the [TopicIndex] this package consumes.
//   - [internal/entity]                    — [Entry],
//     [TopicData], [KeyFileData] domain types.
package moc
