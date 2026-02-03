# Context Health Analysis

## Overview

Analyze Claude Code session transcripts to detect context issues:
- High token usage (approaching model limits)
- Repetition (context degradation signal)

These are **programmatic** checks — cheap, fast, continuous.

## Health Metrics

```go
type ContextHealth struct {
    SessionID     string
    Transcript    string

    // Token metrics
    InputTokens   int       // From most recent assistant message
    OutputTokens  int       // Cumulative
    ContextPct    float64   // Estimated % of model limit

    // Activity
    TurnCount     int
    LastActivity  time.Time
    Duration      time.Duration

    // Quality signals
    RepetitionScore float64  // 0.0 = unique, 1.0 = all identical
}
```

## Token Estimation

Claude Code transcripts include usage data:

```json
{
  "message": {
    "usage": {
      "input_tokens": 150000,
      "output_tokens": 500
    }
  }
}
```

### Approach

Take `input_tokens` from most recent assistant message — this reflects
current context window size.

```go
func EstimateContextPercent(inputTokens int) float64 {
    // Conservative estimate across models
    modelLimit := 200000
    return float64(inputTokens) / float64(modelLimit) * 100
}
```

### Thresholds (with Hysteresis)

| Level | Alert At | Clear At | Code |
|-------|----------|----------|------|
| Info | 70% | 60% | `CTX_HEALTH_70` |
| Warning | 85% | 75% | `CTX_HEALTH_85` |
| Critical | 95% | 85% | `CTX_HEALTH_95` |

## Repetition Detection

Context degradation manifests as repetition — identical responses, tool loops,
apologize cycles.

### Layered Approach

**Level 1 (v1): Exact detection** — cheap, high confidence

```go
func checkExactRepetition(messages []Message, window int) *RepetitionSignal {
    if len(messages) < window {
        return nil
    }

    recent := messages[len(messages)-window:]
    hashes := make(map[uint64][]int)  // hash → turn indices

    for i, msg := range recent {
        h := hashContent(msg.Content)
        if indices, exists := hashes[h]; exists {
            return &RepetitionSignal{
                Type:     "exact_repeat",
                Score:    1.0,
                Evidence: fmt.Sprintf("Turn %d identical to turn %d", i, indices[0]),
            }
        }
        hashes[h] = append(hashes[h], i)
    }

    // Check tool call repetition
    toolCalls := make(map[string]int)  // "tool:args_hash" → count
    for _, msg := range recent {
        for _, call := range msg.ToolCalls {
            key := call.Name + ":" + hashArgs(call.Args)
            toolCalls[key]++
            if toolCalls[key] >= 3 {
                return &RepetitionSignal{
                    Type:     "tool_loop",
                    Score:    0.8,
                    Evidence: fmt.Sprintf("Tool %s called 3+ times with same args", call.Name),
                }
            }
        }
    }

    return nil
}
```

**Level 2 (future): Statistical detection** — for paraphrased repetition

```go
// Future: add when Level 1 proves insufficient
func checkStatisticalRepetition(messages []Message, window int) *RepetitionSignal {
    // N-gram overlap ratio
    // Novel token ratio (unique/total)
    // Hashed TF-IDF cosine similarity
    return nil
}
```

### Detection Targets

| Pattern | Method | Confidence |
|---------|--------|------------|
| Identical responses | Content hash | High |
| Tool loops | Tool+args fingerprint | High |
| Apologize cycles | Phrase detection | Medium |
| Paraphrased repetition | N-gram/TF-IDF (v2) | Medium |

## Transcript Parsing

Stream parse for efficiency:

```go
func AnalyzeTranscript(path string) (*ContextHealth, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    health := &ContextHealth{Transcript: path}
    var messages []Message
    var firstTime time.Time

    scanner := bufio.NewScanner(file)
    // Handle large lines
    buf := make([]byte, 0, 1024*1024)
    scanner.Buffer(buf, 10*1024*1024)

    for scanner.Scan() {
        var msg Message
        if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
            continue
        }

        messages = append(messages, msg)
        health.TurnCount++

        if firstTime.IsZero() {
            firstTime = msg.Timestamp
        }
        health.LastActivity = msg.Timestamp

        if msg.Message.Usage.InputTokens > 0 {
            health.InputTokens = msg.Message.Usage.InputTokens
            health.OutputTokens += msg.Message.Usage.OutputTokens
        }
    }

    health.Duration = health.LastActivity.Sub(firstTime)
    health.ContextPct = EstimateContextPercent(health.InputTokens)

    // Repetition detection (last 10 turns)
    if rep := checkExactRepetition(messages, 10); rep != nil {
        health.RepetitionScore = rep.Score
    }

    return health, nil
}
```

## Alert Generation

```go
func GenerateHealthAlerts(session Session, health *ContextHealth, state *AlertState) []Alert {
    var alerts []Alert

    // Context thresholds (with hysteresis)
    for _, t := range Thresholds {
        key := session.ID + ":" + t.Code
        wasActive := state.ActiveAlerts[key]

        if health.ContextPct >= t.AlertAt && !wasActive {
            // Crossed up — alert
            state.ActiveAlerts[key] = true
            alerts = append(alerts, Alert{
                Auditor:  "context-health",
                Code:     t.Code,
                Session:  session.ID,
                Severity: t.Severity,
                Summary:  fmt.Sprintf("Context: %.0f%% full", health.ContextPct),
                Action:   t.Action,
            })
        } else if health.ContextPct < t.ClearAt && wasActive {
            // Crossed down — clear
            state.ActiveAlerts[key] = false
        }
    }

    // Repetition
    if health.RepetitionScore > 0.5 {
        alerts = append(alerts, Alert{
            Auditor:  "context-health",
            Code:     "CTX_REPETITION",
            Session:  session.ID,
            Severity: "warning",
            Summary:  fmt.Sprintf("Repetition detected (score: %.1f)", health.RepetitionScore),
            Action:   "Check if stuck, try different approach or /clear",
        })
    }

    return alerts
}

var Thresholds = []struct {
    AlertAt  float64
    ClearAt  float64
    Code     string
    Severity string
    Action   string
}{
    {70, 60, "CTX_HEALTH_70", "info", "Monitor usage, consider summarizing soon"},
    {85, 75, "CTX_HEALTH_85", "warning", "Summarize current state, prepare to /clear"},
    {95, 85, "CTX_HEALTH_95", "critical", "/clear now, reload from DECISIONS.md + TASKS.md"},
}
```

## Performance

- **Stream parse**: Don't load entire JSONL
- **Cache by mtime**: Skip unchanged transcripts
- **Window limit**: Repetition on last 10 messages only
- **Early exit**: Stop on first exact match
