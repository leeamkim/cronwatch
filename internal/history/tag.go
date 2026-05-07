package history

import (
	"fmt"
	"sort"
)

// TagFilter returns all entries for a given job that match any of the provided tags.
func TagFilter(store map[string][]Entry, tags []string) ([]Entry, error) {
	if len(tags) == 0 {
		return nil, fmt.Errorf("at least one tag is required")
	}

	tagSet := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		tagSet[t] = struct{}{}
	}

	var results []Entry
	for _, entries := range store {
		for _, e := range entries {
			if entryMatchesTags(e, tagSet) {
				results = append(results, e)
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].StartedAt.Before(results[j].StartedAt)
	})

	return results, nil
}

// TagSummary returns a map of tag -> count of entries carrying that tag.
func TagSummary(store map[string][]Entry) map[string]int {
	summary := make(map[string]int)
	for _, entries := range store {
		for _, e := range entries {
			for _, t := range e.Tags {
				summary[t]++
			}
		}
	}
	return summary
}

func entryMatchesTags(e Entry, tagSet map[string]struct{}) bool {
	for _, t := range e.Tags {
		if _, ok := tagSet[t]; ok {
			return true
		}
	}
	return false
}
