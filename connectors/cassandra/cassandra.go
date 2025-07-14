//
//  cassandra.go
//  coreruntime_extenstions
//
//  Created by karim-w on 14/07/2025.
//

package cassandra

import (
	"crypto/tls"
	"log"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	"github.com/karim-w/coreruntime/coreutils"
)

func connect(
	host string,
	port int,
	keyspace string,
	username string,
	password string,
	tlsEnabled bool,
) (*gocql.Session, error) {
	// Connect to Cassandra cluster:
	cluster := gocql.NewCluster(
		coreutils.StringBuilder(host, ":", strconv.Itoa(port)),
	)
	// cluster.Port = port

	log.Printf(
		"Connecting to Cassandra cluster... host:%s, port:%d, username:%s, keyspace: %s\n",
		host,
		port,
		username,
		keyspace,
	)

	// Set cluster config:
	// cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 30
	cluster.Keyspace = keyspace

	if tlsEnabled {
		log.Println("Enabling TLS for Cassandra connection...")
		cluster.SslOpts = &gocql.SslOptions{
			Config: &tls.Config{
				MinVersion:         tls.VersionTLS12,
				InsecureSkipVerify: true,
			},
		}
	}

	// cluster.CQLVersion = "3.11.0"
	// cluster.DisableInitialHostLookup = true

	// Set authentication:
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}

	// Set port:

	// Create session:
	session, err := cluster.CreateSession()
	if err != nil {
		log.Println("Error creating session", err)
		return nil, err
	}

	log.Println("Connected to Cassandra cluster")

	return session, nil
}
