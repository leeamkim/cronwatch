# cronwatch

Lightweight cron job monitor that sends alerts when jobs fail or run too long.

## Installation

```bash
go install github.com/yourname/cronwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/cronwatch.git && cd cronwatch && go build ./...
```

## Usage

Wrap your cron job command with `cronwatch` to start monitoring it:

```bash
cronwatch --name "daily-backup" --timeout 30m -- /usr/local/bin/backup.sh
```

Configure alerts in `cronwatch.yaml`:

```yaml
alerts:
  slack:
    webhook_url: "https://hooks.slack.com/services/..."
  email:
    to: "ops@example.com"

jobs:
  daily-backup:
    timeout: 30m
    notify_on: [failure, timeout]
```

Then run your cron job as usual:

```
0 2 * * * cronwatch --name "daily-backup" -- /usr/local/bin/backup.sh
```

cronwatch will send an alert if the job exits with a non-zero status or exceeds the configured timeout.

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--name` | Job identifier | required |
| `--timeout` | Max allowed runtime | `0` (disabled) |
| `--config` | Path to config file | `./cronwatch.yaml` |

## License

MIT © [yourname](https://github.com/yourname)