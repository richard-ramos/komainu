package sqlite

import (
	"database/sql"
	"strings"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlcipher"
)

func Driver(db *sql.DB) (database.Driver, error) {
	return sqlcipher.WithInstance(db, &sqlcipher.Config{
		MigrationsTable: sqlcipher.DefaultMigrationsTable,
	})
}

func URLDefaults(dburl string) string {
	if !strings.Contains(dburl, "?") {
		dburl += "?"
	}

	if !strings.Contains(dburl, "_journal=") {
		dburl += "&_journal=WAL"
	}

	if !strings.Contains(dburl, "_timeout=") {
		dburl += "&_timeout=5000"
	}

	return dburl
}
