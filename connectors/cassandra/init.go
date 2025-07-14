//
//  init.go
//  coreruntime_extenstions
//
//  Created by karim-w on 14/07/2025.
//

package cassandra

import (
	"context"
	"log/slog"
	"time"

	"github.com/gocql/gocql"
	"github.com/karim-w/coreruntime_extensions/common/options"
)

// Initialize establishes a connection to a Cassandra database using the
// provided connection string and options.
// The connection string should be in the format:
// cassandra://host:port/keyspace?username=username&password=password&ssl=true
// there are options for panicable pings and other configurations.
func Initialize(
	connection_string string,
	options options.Cassandra,
) (*gocql.Session, error) {
	// Parse the connection string to extract the necessary parameters
	username, password, host, keyspace, port, tlsEnabled, err := parseUri(connection_string)
	if err != nil {
		return nil, err
	}

	session, err := connect(host, port, keyspace, username, password, tlsEnabled)
	if err != nil {
		return nil, err
	}

	if err := session.Query("SELECT now() FROM system.local").Exec(); err != nil {
		return nil, err
	}

	if options.PanicablePings {
		go poll_db_heath(session)
	}

	return session, nil
}

func poll_db_heath(cql *gocql.Session) {
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
			if err := cql.Query("SELECT now() FROM system.local").Exec(); err != nil {
				slog.Error("db ping failed", "error", err)
				continue
			}
			// successful ping: reset timer
			timer.Reset(time.Second * 15)
		}
	}
}
