package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/user/cronwatch/internal/history"
	"github.com/urfave/cli/v2"
)

var tagCmd = &cli.Command{
	Name:  "tag",
	Usage: "filter or summarise history entries by tag",
	Subcommands: []*cli.Command{
		{
			Name:      "filter",
			Usage:     "list entries matching one or more tags",
			ArgsUsage: "<tag> [tag...]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "history",
					Value: defaultHistoryPath(),
					Usage: "path to history file",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					return fmt.Errorf("at least one tag argument is required")
				}
				h, err := history.New(c.String("history"))
				if err != nil {
					return err
				}
				entries, err := history.TagFilter(h.Store(), c.Args().Slice())
				if err != nil {
					return err
				}
				if len(entries) == 0 {
					fmt.Println("no entries found for the given tags")
					return nil
				}
				fmt.Fprintf(os.Stdout, "%-20s %-25s %-8s %s\n", "JOB", "STARTED", "STATUS", "TAGS")
				for _, e := range entries {
					status := "ok"
					if !e.Success {
						status = "fail"
					}
					fmt.Fprintf(os.Stdout, "%-20s %-25s %-8s %s\n",
						e.Job,
						e.StartedAt.Format("2006-01-02 15:04:05"),
						status,
						strings.Join(e.Tags, ","),
					)
				}
				return nil
			},
		},
		{
			Name:  "summary",
			Usage: "show count of entries per tag",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "history",
					Value: defaultHistoryPath(),
					Usage: "path to history file",
				},
			},
			Action: func(c *cli.Context) error {
				h, err := history.New(c.String("history"))
				if err != nil {
					return err
				}
				summary := history.TagSummary(h.Store())
				if len(summary) == 0 {
					fmt.Println("no tags found")
					return nil
				}
				keys := make([]string, 0, len(summary))
				for k := range summary {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				fmt.Fprintf(os.Stdout, "%-20s %s\n", "TAG", "COUNT")
				for _, k := range keys {
					fmt.Fprintf(os.Stdout, "%-20s %d\n", k, summary[k])
				}
				return nil
			},
		},
	},
}
