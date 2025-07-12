//
//  psql.go
//  coreruntime_extenstions
//
//  Created by karim-w on 12/07/2025.
//

package psql

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/karim-w/coreruntime_extensions/common/options"
	_ "github.com/lib/pq"
)

const (
	DEFAULT_MAX_IDLE_CONNS = 10
	DEFAULT_MAX_OPEN_CONNS = 100
)

func MustInit(
	dsn string,
	opt options.SQL_Database,
) (*sql.DB, error) {
	idle_conns := opt.MaxIdleConns.GetOrElse(DEFAULT_MAX_IDLE_CONNS)
	open_conns := opt.MaxOpenConns.GetOrElse(DEFAULT_MAX_OPEN_CONNS)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(idle_conns)
	db.SetMaxOpenConns(open_conns)

	err = db.Ping()
	if err != nil {
		slog.Error("failed to ping database", "error", err)
		return nil, err
	}

	if opt.PanicablePings {
		go poll_db_heath(db)
	}

	return db, nil
}

func poll_db_heath(db *sql.DB) {
	// This context is canceled only when the timeout expires without a successful ping
	ctx, cancel := context.WithCancel(context.Background())

	timer := time.AfterFunc(time.Second*15, func() {
		// Triggered only if no successful ping resets the timer in time
		panic("db unresponsive for too long")
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
			if err := db.Ping(); err != nil {
				slog.Error("db ping failed", "error", err)
				// optional: log the error here
				continue
			}
			// successful ping: reset timer
			timer.Reset(time.Second * 15)
		}
	}
}
