//
//  rdb.go
//  coreruntime_extenstions
//
//  Created by karim-w on 12/07/2025.
//

package rdb

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/karim-w/coreruntime_extensions/common/options"
)

func Initialize(
	url string,
	opt options.Redis,
) (*redis.Client, error) {
	cfg, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(cfg)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	if opt.PanicablePings {
		go poll_cache_heath(client)
	}

	return client, nil
}

func poll_cache_heath(rdb *redis.Client) {
	// This context is canceled only when the timeout expires without a successful ping
	ctx, cancel := context.WithCancel(context.Background())

	timer := time.AfterFunc(time.Second*15, func() {
		// Triggered only if no successful ping resets the timer in time
		panic("cache unresponsive for too long")
	})

	defer func() {
		timer.Stop()
		cancel()
	}()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := rdb.Ping(ctx).Err(); err != nil {
				slog.Error("cache ping failed", "error", err)
				// optional: log the error here
				continue
			}
			// successful ping: reset timer
			timer.Reset(time.Second * 15)
		}
	}
}
