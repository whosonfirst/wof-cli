package repos

import (
	"context"
	"fmt"
	_ "log/slog"
	"regexp"
	"strconv"
	"time"

	"github.com/sfomuseum/iso8601duration"
	"github.com/whosonfirst/go-whosonfirst-github/organizations"
	"github.com/whosonfirst/wof"
)

type ReposCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "repos", NewReposCommand)
}

func NewReposCommand(ctx context.Context, cmd string) (wof.Command, error) {

	c := &ReposCommand{}
	return c, nil
}

func (c *ReposCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	opts := organizations.NewDefaultListOptions()

	opts.Prefix = prefix
	opts.Exclude = exclude
	opts.Forked = forked
	opts.NotForked = not_forked
	opts.AccessToken = token
	opts.Debug = verbose // This actually means verbose (and this package should be updated accordingly...)
	opts.EnsureCommits = ensure_commits
	opts.ExcludeArchived = exclude_archived

	if updated_since != "" {

		var since time.Time

		is_timestamp, err := regexp.MatchString("^\\d+$", updated_since)

		if err != nil {
			return fmt.Errorf("Updated since regexp failed, %w", err)
		}

		if is_timestamp {

			ts, err := strconv.Atoi(updated_since)

			if err != nil {
				return fmt.Errorf("Failed to parse updated_since value, %w", err)
			}

			now := time.Now()

			tm := time.Unix(int64(ts), 0)
			since = now.Add(-time.Since(tm))

		} else {

			// maybe also this https://github.com/araddon/dateparse ?

			d, err := duration.FromString(updated_since)

			if err != nil {
				return fmt.Errorf("Failed to parse updated since value, %w", err)
			}

			now := time.Now()
			since = now.Add(-d.ToDuration())
		}

		opts.PushedSince = &since
	}

	repos, err := organizations.ListRepos(org, opts)

	if err != nil {
		fmt.Errorf("Failed to list repos, %w", err)
	}

	for _, name := range repos {
		fmt.Println(name)
	}

	return nil
}
