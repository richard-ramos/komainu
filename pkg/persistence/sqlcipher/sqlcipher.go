package sqlcipher

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlcipher"
)

func Driver(db *sql.DB) (database.Driver, error) {
	return sqlcipher.WithInstance(db, &sqlcipher.Config{
		MigrationsTable: sqlcipher.DefaultMigrationsTable,
	})
}

func URLDefaults(dburl string, key string, kdfIterations int) string {
	if !strings.Contains(dburl, "?") {
		dburl += "?"
	}

	if !strings.Contains(dburl, "_pragma_key=") {
		dburl += "&_pragma_key=" + key
	}

	if !strings.Contains(dburl, "_kdf_iter=") {
		dburl += fmt.Sprintf("&_kdf_iter=%d", kdfIterations)
	}

	if !strings.Contains(dburl, "_journal=") {
		dburl += "&_journal=WAL"
	}

	if !strings.Contains(dburl, "_timeout=") {
		dburl += "&_timeout=5000"
	}

	return dburl
}
