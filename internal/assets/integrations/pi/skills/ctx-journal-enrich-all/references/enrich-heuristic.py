#!/usr/bin/env python3
# Copyright 2026 ActiveMemory. All rights reserved.
#
# Heuristic journal enrichment script.
# Adds type/outcome/topics/technologies/summary frontmatter fields
# to journal entries based on title and filename pattern matching.
#
# Usage:
#   python3 enrich-heuristic.py <file-list.txt>
#
# The file list should contain one journal file path per line.
# Files already containing type: and outcome: fields are skipped.
#
# After enrichment, each file is marked via:
#   ctx system mark-journal <filename> enriched

import os
import re
import subprocess
import sys

# --- Detection heuristics ---

TYPE_KEYWORDS = [
    (["fix", "bug", "broken", "debug", "crash", "oom", "error"], "bugfix"),
    (["implement", "add", "create", "build", "absorb", "convert"], "feature"),
    (["refactor", "rename", "restructure", "reorganize", "consolidate", "migrate", "move"], "refactor"),
    (["plan", "design", "spec", "brainstorm", "explore", "investigate", "evaluate", "triage"], "planning"),
    (["doc", "recipe", "blog", "contributing", "navigation", "admonition", "changelog", "publish"], "documentation"),
    (["audit", "review", "verify", "check", "sanitize", "lint", "batch", "archive", "enrich"], "maintenance"),
]

OUTCOME_FILENAME_PATTERNS = {
    "request-interrupted": "abandoned",
    "clear-clear": "abandoned",
    "brief-session": "partial",
}

TOPIC_KEYWORDS = [
    (["journal", "enrich"], "journal"),
    (["task", "archive"], "task-management"),
    (["hook", "nudge"], "hooks"),
    (["skill"], "skills"),
    (["doc", "recipe", "blog"], "documentation"),
    (["recall", "session", "remember"], "session-history"),
    (["encrypt", "key", "pad"], "encryption"),
    (["drift"], "context-drift"),
    (["webhook", "notify"], "notifications"),
    (["worktree"], "git-worktrees"),
    (["permission", "sanitize"], "security"),
    (["site", "render", "nav"], "site-generation"),
    (["init", "bootstrap"], "initialization"),
    (["config", "ctxrc"], "configuration"),
    (["test", "lint"], "testing"),
    (["export", "import"], "data-pipeline"),
    (["lock", "unlock"], "data-safety"),
    (["remind"], "reminders"),
    (["resource", "oom"], "system-resources"),
    (["map", "architecture"], "architecture"),
    (["plan", "spec"], "planning"),
    (["commit", "release"], "version-control"),
    (["prompt"], "prompt-templates"),
    (["rss", "feed"], "rss-feed"),
    (["model", "opus"], "ai-models"),
    (["context"], "context-management"),
    (["cli", "command"], "cli"),
    (["hack", "script", "absorb"], "build-tooling"),
]

TECH_KEYWORDS = [
    (["bash", "shell", "script", "hack"], "bash"),
    (["yaml", "ctxrc"], "yaml"),
    (["json", "schema"], "json"),
    (["markdown", "mkdocs", "zensical", "site", "render"], "markdown"),
    (["git", "commit", "worktree"], "git"),
    (["webhook", "http"], "http"),
    (["rss", "atom", "feed"], "rss"),
    (["encrypt", "aes", "key", "crypto"], "aes-256-gcm"),
]


def get_title(fm_block):
    m = re.search(r'^title:\s*"?(.+?)"?\s*$', fm_block, re.M)
    return m.group(1) if m else ""


def has_enrichment(fm_block):
    return bool(re.search(r'^type:', fm_block, re.M)) and bool(
        re.search(r'^outcome:', fm_block, re.M)
    )


def detect_type(title):
    t = title.lower()
    for keywords, typ in TYPE_KEYWORDS:
        if any(w in t for w in keywords):
            return typ
    return "exploration"


def detect_outcome(filename):
    fname = filename.lower()
    for pattern, outcome in OUTCOME_FILENAME_PATTERNS.items():
        if pattern in fname:
            return outcome
    return "completed"


def detect_topics(title):
    t = title.lower()
    topics = []
    for keywords, topic in TOPIC_KEYWORDS:
        if any(k in t for k in keywords):
            topics.append(topic)
    return topics[:5] if topics else ["general"]


def detect_technologies(title):
    t = title.lower()
    techs = {"go", "cli"}  # ctx defaults
    for keywords, tech in TECH_KEYWORDS:
        if any(w in t for w in keywords):
            techs.add(tech)
    return sorted(techs)


def enrich_file(filepath):
    with open(filepath) as f:
        content = f.read()

    if not content.startswith("---"):
        return False

    # Split on closing --- of frontmatter (first \n--- after opening ---)
    try:
        idx = content.index("\n---", 3)
    except ValueError:
        return False

    fm_block = content[3:idx]
    rest = content[idx:]

    if has_enrichment(fm_block):
        return False

    title = get_title(fm_block)
    fname = os.path.basename(filepath)

    typ = detect_type(title)
    outcome = detect_outcome(fname)
    topics = detect_topics(title)
    techs = detect_technologies(title)

    fields = f"type: {typ}\noutcome: {outcome}\n"
    fields += "topics:\n" + "".join(f"  - {t}\n" for t in topics)
    fields += "technologies:\n" + "".join(f"  - {t}\n" for t in techs)
    fields += f'summary: "{title}"\n'

    new_content = "---" + fm_block + "\n" + fields + rest

    with open(filepath, "w") as f:
        f.write(new_content)

    subprocess.run(
        ["ctx", "system", "mark-journal", fname, "enriched"],
        capture_output=True,
    )
    return True


def main():
    if len(sys.argv) < 2:
        print("Usage: enrich-heuristic.py <file-list.txt>", file=sys.stderr)
        sys.exit(1)

    with open(sys.argv[1]) as f:
        files = [line.strip() for line in f if line.strip()]

    enriched = 0
    skipped = 0
    for filepath in files:
        if not os.path.exists(filepath):
            print(f"SKIP {filepath} (not found)")
            skipped += 1
            continue
        if enrich_file(filepath):
            enriched += 1
            print(f"OK   {os.path.basename(filepath)}")
        else:
            skipped += 1

    print(f"\nEnriched: {enriched}, Skipped: {skipped}, Total: {len(files)}")


if __name__ == "__main__":
    main()
