package history

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// AlertRule defines when an alert threshold is breached.
type AlertRule struct {
	JobName         string
	MaxFailStreak   int
	MinSuccessRate  float64 // 0–1
	MaxAvgDuration  time.Duration
}

// AlertBreach describes a single threshold violation.
type AlertBreach struct {
	JobName string
	Reason  string
	At      time.Time
}

// CheckAlerts evaluates AlertRules against stored history and returns any
// breaches. path is the history JSON file.
func CheckAlerts(path string, rules []AlertRule) ([]AlertBreach, error) {
	if path == "" {
		return nil, fmt.Errorf("history path must not be empty")
	}
	if len(rules) == 0 {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read history: %w", err)
	}

	var store map[string][]Entry
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("parse history: %w", err)
	}

	var breaches []AlertBreach
	for _, rule := range rules {
		entries := store[rule.JobName]
		if len(entries) == 0 {
			continue
		}

		if rule.MaxFailStreak > 0 {
			if streak := trailingFailStreak(entries); streak >= rule.MaxFailStreak {
				breaches = append(breaches, AlertBreach{
					JobName: rule.JobName,
					Reason:  fmt.Sprintf("failure streak of %d reached threshold %d", streak, rule.MaxFailStreak),
					At:      time.Now(),
				})
			}
		}

		if rule.MinSuccessRate > 0 {
			rate := successRate(entries)
			if rate < rule.MinSuccessRate {
				breaches = append(breaches, AlertBreach{
					JobName: rule.JobName,
					Reason:  fmt.Sprintf("success rate %.0f%% below threshold %.0f%%", rate*100, rule.MinSuccessRate*100),
					At:      time.Now(),
				})
			}
		}

		if rule.MaxAvgDuration > 0 {
			if avg := avgDuration(entries); avg > rule.MaxAvgDuration {
				breaches = append(breaches, AlertBreach{
					JobName: rule.JobName,
					Reason:  fmt.Sprintf("avg duration %s exceeds threshold %s", avg.Round(time.Millisecond), rule.MaxAvgDuration),
					At:      time.Now(),
				})
			}
		}
	}
	return breaches, nil
}

func trailingFailStreak(entries []Entry) int {
	count := 0
	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].Status != "failure" {
			break
		}
		count++
	}
	return count
}

func successRate(entries []Entry) float64 {
	if len(entries) == 0 {
		return 0
	}
	ok := 0
	for _, e := range entries {
		if e.Status == "success" {
			ok++
		}
	}
	return float64(ok) / float64(len(entries))
}

func avgDuration(entries []Entry) time.Duration {
	if len(entries) == 0 {
		return 0
	}
	var total time.Duration
	for _, e := range entries {
		total += e.Duration
	}
	return total / time.Duration(len(entries))
}
