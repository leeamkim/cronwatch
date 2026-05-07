package history

import (
	"testing"
	"time"
)

func seedStats(t *testing.T) *Store {
	t.Helper()
	s, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	now := time.Now()
	entries := []Entry{
		{JobName: "backup", StartedAt: now.Add(-3 * time.Hour), Duration: 2 * time.Second, Error: ""},
		{JobName: "backup", StartedAt: now.Add(-2 * time.Hour), Duration: 4 * time.Second, Error: "exit 1"},
		{JobName: "backup", StartedAt: now.Add(-1 * time.Hour), Duration: 3 * time.Second, Error: ""},
		{JobName: "sync", StartedAt: now.Add(-30 * time.Minute), Duration: 500 * time.Millisecond, Error: ""},
	}
	for _, e := range entries {
		s.data[e.JobName] = append(s.data[e.JobName], e)
	}
	return s
}

func TestComputeStats_CountsRuns(t *testing.T) {
	s := seedStats(t)
	stats := ComputeStats(s)

	backup, ok := stats["backup"]
	if !ok {
		t.Fatal("expected stats for 'backup'")
	}
	if backup.TotalRuns != 3 {
		t.Errorf("TotalRuns = %d, want 3", backup.TotalRuns)
	}
	if backup.SuccessCount != 2 {
		t.Errorf("SuccessCount = %d, want 2", backup.SuccessCount)
	}
	if backup.FailureCount != 1 {
		t.Errorf("FailureCount = %d, want 1", backup.FailureCount)
	}
}

func TestComputeStats_AvgAndMaxDuration(t *testing.T) {
	s := seedStats(t)
	stats := ComputeStats(s)

	backup := stats["backup"]
	if backup.MaxDuration != 4*time.Second {
		t.Errorf("MaxDuration = %v, want 4s", backup.MaxDuration)
	}
	wantAvg := (2 + 4 + 3) * time.Second / 3
	if backup.AvgDuration != wantAvg {
		t.Errorf("AvgDuration = %v, want %v", backup.AvgDuration, wantAvg)
	}
}

func TestComputeStats_SuccessRate(t *testing.T) {
	s := seedStats(t)
	stats := ComputeStats(s)

	backup := stats["backup"]
	got := backup.SuccessRate()
	want := 66.66666666666667
	if got < 66.6 || got > 66.7 {
		t.Errorf("SuccessRate = %.2f, want ~%.2f", got, want)
	}

	sync := stats["sync"]
	if sync.SuccessRate() != 100.0 {
		t.Errorf("sync SuccessRate = %.2f, want 100", sync.SuccessRate())
	}
}

func TestComputeStats_LastStatus(t *testing.T) {
	s := seedStats(t)
	stats := ComputeStats(s)

	backup := stats["backup"]
	if backup.LastStatus != "ok" {
		t.Errorf("LastStatus = %q, want \"ok\"", backup.LastStatus)
	}
}

func TestJobStats_String_ContainsJobName(t *testing.T) {
	stats := JobStats{
		JobName:      "myjob",
		TotalRuns:    5,
		SuccessCount: 4,
		FailureCount: 1,
		AvgDuration:  time.Second,
		MaxDuration:  2 * time.Second,
		LastRun:      time.Now(),
		LastStatus:   "ok",
	}
	str := stats.String()
	if len(str) == 0 {
		t.Error("expected non-empty string")
	}
	if str[:5] != "myjob" {
		t.Errorf("String() does not start with job name: %s", str)
	}
}
