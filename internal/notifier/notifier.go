package notifier

import (
	"fmt"
	"io"
	"os"
	"time"
)

// AlertType represents the kind of alert being sent.
type AlertType string

const (
	AlertFailure AlertType = "FAILURE"
	AlertTimeout AlertType = "TIMEOUT"
)

// Alert holds the details of a cron job alert.
type Alert struct {
	JobName   string
	AlertType AlertType
	Message   string
	OccuredAt time.Time
}

// Notifier sends alerts through a configured output.
type Notifier struct {
	writer io.Writer
}

// New creates a Notifier that writes to the given writer.
// If writer is nil, os.Stderr is used.
func New(writer io.Writer) *Notifier {
	if writer == nil {
		writer = os.Stderr
	}
	return &Notifier{writer: writer}
}

// Notify formats and sends an alert.
func (n *Notifier) Notify(alert Alert) error {
	if alert.JobName == "" {
		return fmt.Errorf("notifier: job name must not be empty")
	}

	timestamp := alert.OccuredAt
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	line := fmt.Sprintf(
		"[cronwatch] %s | job=%q type=%s message=%q\n",
		timestamp.UTC().Format(time.RFC3339),
		alert.JobName,
		alert.AlertType,
		alert.Message,
	)

	_, err := fmt.Fprint(n.writer, line)
	if err != nil {
		return fmt.Errorf("notifier: failed to write alert: %w", err)
	}
	return nil
}
