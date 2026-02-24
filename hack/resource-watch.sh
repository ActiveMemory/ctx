#!/bin/bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

# Continuous system resource monitor with ASCII bars.
# Uses the same thresholds as ctx system resources.
# Usage: ./hack/resource-watch.sh [interval_seconds]
#
# Thresholds (from ctx sysinfo):
#   Memory:  WARNING >= 80%    DANGER >= 90%
#   Swap:    WARNING >= 50%    DANGER >= 75%
#   Disk:    WARNING >= 85%    DANGER >= 95%
#   Load:    WARNING >= 0.8x   DANGER >= 1.5x CPUs

set -euo pipefail

INTERVAL="${1:-5}"
BAR_WIDTH=30

# Colors
GREEN="\033[32m"
YELLOW="\033[33m"
RED="\033[31m"
DIM="\033[2m"
BOLD="\033[1m"
RESET="\033[0m"

# ── helpers ──────────────────────────────────────────────────────────

bar() {
  local pct=$1 width=$2
  local filled=$((pct * width / 100))
  [ "$filled" -gt "$width" ] && filled=$width
  [ "$filled" -lt 0 ] && filled=0
  local empty=$((width - filled))
  printf '%0.s█' $(seq 1 "$filled" 2>/dev/null)
  printf '%0.s░' $(seq 1 "$empty" 2>/dev/null)
}

color_for() {
  local resource=$1 pct=$2
  case "$resource" in
    mem)  [ "$pct" -ge 90 ] && echo "$RED"  && return
          [ "$pct" -ge 80 ] && echo "$YELLOW" && return ;;
    swap) [ "$pct" -ge 75 ] && echo "$RED"  && return
          [ "$pct" -ge 50 ] && echo "$YELLOW" && return ;;
    disk) [ "$pct" -ge 95 ] && echo "$RED"  && return
          [ "$pct" -ge 85 ] && echo "$YELLOW" && return ;;
  esac
  echo "$GREEN"
}

icon_for() {
  local resource=$1 pct=$2
  case "$resource" in
    mem)  [ "$pct" -ge 90 ] && echo -e "${RED}✗${RESET}" && return
          [ "$pct" -ge 80 ] && echo -e "${YELLOW}⚠${RESET}" && return ;;
    swap) [ "$pct" -ge 75 ] && echo -e "${RED}✗${RESET}" && return
          [ "$pct" -ge 50 ] && echo -e "${YELLOW}⚠${RESET}" && return ;;
    disk) [ "$pct" -ge 95 ] && echo -e "${RED}✗${RESET}" && return
          [ "$pct" -ge 85 ] && echo -e "${YELLOW}⚠${RESET}" && return ;;
  esac
  echo -e "${GREEN}✓${RESET}"
}

fmt_gib() {
  local kb=$1
  awk "BEGIN { printf \"%.1f\", $kb / 1048576 }"
}

# ── readers ──────────────────────────────────────────────────────────

read_memory() {
  local key val _
  MEM_TOTAL=0; MEM_AVAIL=0; SWAP_TOTAL=0; SWAP_FREE=0
  while read -r key val _; do
    case "$key" in
      MemTotal:)     MEM_TOTAL=$val ;;
      MemAvailable:) MEM_AVAIL=$val ;;
      SwapTotal:)    SWAP_TOTAL=$val ;;
      SwapFree:)     SWAP_FREE=$val ;;
    esac
  done < /proc/meminfo
  MEM_USED=$((MEM_TOTAL - MEM_AVAIL))
  SWAP_USED=$((SWAP_TOTAL - SWAP_FREE))
}

read_load() {
  read -r LOAD1 LOAD5 LOAD15 _ < /proc/loadavg
  NUM_CPU=$(nproc 2>/dev/null || getconf _NPROCESSORS_ONLN 2>/dev/null || echo 1)
}

read_disk() {
  local _ total used avail
  read -r _ total used avail _ < <(df -k . | tail -1)
  DISK_TOTAL=$total
  DISK_USED=$used
}

read_uptime() {
  local raw
  read -r raw _ < /proc/uptime
  local secs=${raw%%.*}
  local d=$((secs / 86400))
  local h=$(( (secs % 86400) / 3600 ))
  local m=$(( (secs % 3600) / 60 ))
  UPTIME="${d}d ${h}h ${m}m"
}

# ── main loop ────────────────────────────────────────────────────────

while true; do
  clear

  read_memory
  read_load
  read_disk
  read_uptime

  # Percentages
  MEM_PCT=0; SWAP_PCT=0; DISK_PCT=0
  [ "$MEM_TOTAL"  -gt 0 ] && MEM_PCT=$((MEM_USED * 100 / MEM_TOTAL))
  [ "$SWAP_TOTAL" -gt 0 ] && SWAP_PCT=$((SWAP_USED * 100 / SWAP_TOTAL))
  [ "$DISK_TOTAL" -gt 0 ] && DISK_PCT=$((DISK_USED * 100 / DISK_TOTAL))

  # Load ratio (integer × 100 for comparison)
  LOAD_RATIO_100=$(awk "BEGIN { printf \"%d\", $LOAD1 / $NUM_CPU * 100 }")
  LOAD_RATIO=$(awk "BEGIN { printf \"%.2f\", $LOAD1 / $NUM_CPU }")
  LOAD_PCT=$((LOAD_RATIO_100 * 100 / 150))  # scale: 1.5x = 100%
  [ "$LOAD_PCT" -gt 100 ] && LOAD_PCT=100

  # Overall status
  WORST="ok"
  for check in "$MEM_PCT:90:80" "$SWAP_PCT:75:50" "$DISK_PCT:95:85"; do
    IFS=: read -r val danger warn <<< "$check"
    [ "$val" -ge "$danger" ] && WORST="DANGER"
    [ "$val" -ge "$warn" ] && [ "$WORST" != "DANGER" ] && WORST="WARNING"
  done
  if [ "$LOAD_RATIO_100" -ge 150 ]; then
    WORST="DANGER"
  elif [ "$LOAD_RATIO_100" -ge 80 ] && [ "$WORST" != "DANGER" ]; then
    WORST="WARNING"
  fi

  case "$WORST" in
    DANGER)  STATUS_COLOR="$RED" ;;
    WARNING) STATUS_COLOR="$YELLOW" ;;
    *)       STATUS_COLOR="$GREEN"; WORST="ALL CLEAR" ;;
  esac

  # Load color
  LOAD_COLOR="$GREEN"
  LOAD_ICON="${GREEN}✓${RESET}"
  if [ "$LOAD_RATIO_100" -ge 150 ]; then
    LOAD_COLOR="$RED"; LOAD_ICON="${RED}✗${RESET}"
  elif [ "$LOAD_RATIO_100" -ge 80 ]; then
    LOAD_COLOR="$YELLOW"; LOAD_ICON="${YELLOW}⚠${RESET}"
  fi

  # Format values
  MEM_USED_G=$(fmt_gib "$MEM_USED")
  MEM_TOTAL_G=$(fmt_gib "$MEM_TOTAL")
  SWAP_USED_G=$(fmt_gib "$SWAP_USED")
  SWAP_TOTAL_G=$(fmt_gib "$SWAP_TOTAL")
  DISK_USED_G=$(fmt_gib "$DISK_USED")
  DISK_TOTAL_G=$(fmt_gib "$DISK_TOTAL")

  # ── render ───────────────────────────────────────────────────────

  echo ""
  echo -e "  ${BOLD}Resource Monitor${RESET}  ${STATUS_COLOR}[$WORST]${RESET}"
  echo ""

  MC=$(color_for mem  "$MEM_PCT")
  SC=$(color_for swap "$SWAP_PCT")
  DC=$(color_for disk "$DISK_PCT")
  MI=$(icon_for mem  "$MEM_PCT")
  SI=$(icon_for swap "$SWAP_PCT")
  DI=$(icon_for disk "$DISK_PCT")

  printf "  MEM   %b  %5s / %5s GB  (%2d%%)  %b\n" \
    "${MC}$(bar "$MEM_PCT"  "$BAR_WIDTH")${RESET}" \
    "$MEM_USED_G" "$MEM_TOTAL_G" "$MEM_PCT" "$MI"

  printf "  SWAP  %b  %5s / %5s GB  (%2d%%)  %b\n" \
    "${SC}$(bar "$SWAP_PCT" "$BAR_WIDTH")${RESET}" \
    "$SWAP_USED_G" "$SWAP_TOTAL_G" "$SWAP_PCT" "$SI"

  printf "  DISK  %b  %5s / %5s GB  (%2d%%)  %b\n" \
    "${DC}$(bar "$DISK_PCT" "$BAR_WIDTH")${RESET}" \
    "$DISK_USED_G" "$DISK_TOTAL_G" "$DISK_PCT" "$DI"

  printf "  LOAD  %b  %sx  (%d CPUs)  %b\n" \
    "${LOAD_COLOR}$(bar "$LOAD_PCT" "$BAR_WIDTH")${RESET}" \
    "$LOAD_RATIO" "$NUM_CPU" "$LOAD_ICON"

  echo ""
  echo -e "  ${DIM}Load avg:  $LOAD1 / $LOAD5 / $LOAD15  (1m / 5m / 15m)${RESET}"
  echo -e "  ${DIM}Uptime:    $UPTIME${RESET}"
  echo ""
  echo -e "  Refreshing every ${INTERVAL}s. Ctrl+C to stop."

  sleep "$INTERVAL"
done
