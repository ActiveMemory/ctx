# Load Alert: Use 5-Minute Average

## Problem

The resource load alert uses the 1-minute load average, which spikes
during normal build/test runs and produces false positives. A single
`make test` run triggers "critically low resources" even when the
system is healthy.

## Approach

Switch from `Load1` to `Load5` in `sysinfo.Evaluate`. The 5-minute
average smooths transient build/test spikes while still catching
sustained resource pressure.

## Non-Goals

- Changing the threshold ratios (0.8 warn, 1.5 danger)
- Adding configurable averaging windows
