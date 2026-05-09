package history

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// alertRuleYAML mirrors AlertRule for YAML unmarshalling.
type alertRuleYAML struct {
	JobName        string  `yaml:"job"`
	MaxFailStreak  int     `yaml:"max_fail_streak"`
	MinSuccessRate float64 `yaml:"min_success_rate"`
	MaxAvgDuration string  `yaml:"max_avg_duration"`
}

type alertConfigFile struct {
	Rules []alertRuleYAML `yaml:"alert_rules"`
}

// LoadAlertRules reads alert rules from a YAML config file.
func LoadAlertRules(path string) ([]AlertRule, error) {
	if path == "" {
		return nil, fmt.Errorf("alert config path must not be empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read alert config: %w", err)
	}

	var cfg alertConfigFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse alert config: %w", err)
	}

	var rules []AlertRule
	for _, r := range cfg.Rules {
		if r.JobName == "" {
			return nil, fmt.Errorf("alert rule missing job name")
		}
		rule := AlertRule{
			JobName:        r.JobName,
			MaxFailStreak:  r.MaxFailStreak,
			MinSuccessRate: r.MinSuccessRate,
		}
		if r.MaxAvgDuration != "" {
			d, err := time.ParseDuration(r.MaxAvgDuration)
			if err != nil {
				return nil, fmt.Errorf("invalid max_avg_duration for job %q: %w", r.JobName, err)
			}
			rule.MaxAvgDuration = d
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
