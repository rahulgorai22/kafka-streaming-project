package utils

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

// SensorData represents your data structure
type SensorData struct {
	SHA256             string    `db:"sha256"`
	MaliciousnessScore float64   `db:"maliciousness_score"`
	Classification     string    `db:"classification"`
	Timestamp          time.Time `db:"ts"`
}

type CassandraUtilities struct {
	// Define additional cassandra utilities
}

var (
	cassandraHosts = flag.String("cassandra-hosts", "cassandra", "list of cassandra hosts to connect to")
)

func (cu *CassandraUtilities) CreateDatabaseSession() (*gocql.Session, error) {
	cluster := gocql.NewCluster(strings.Split(*cassandraHosts, ",")...)
	cluster.Keyspace = DATABASE_KEYSPACE
	session, err := cluster.CreateSession()
	return session, err
}

func (cu *CassandraUtilities) CloseSession(session *gocql.Session) {
	session.Close()
}

func (cu *CassandraUtilities) InsertData(session *gocql.Session, sha256 string, score float64, classification string, timestamp string) error {

	// Convert the score to float32 as float32 is supported in cassandra.
	score32 := float32(score)
	// Parse the timestamp string into time.Time
	ts, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %v", err)
	}

	// Convert the timestamp to Unix milliseconds.
	timestampMillis := ts.UnixNano() / int64(time.Millisecond)

	query, names := qb.Insert(CLASSIFICATION_TABLE).Columns(SHA256, MALICIOUS_SCORE, CLASSIFICATION, TIMESTAMP).ToCql()
	err = gocqlx.Query(session.Query(query, sha256, score32, classification, timestampMillis), names).ExecRelease()
	return err
}
