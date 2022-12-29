package sqlcipher

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlcipher"
)

func Driver(tablePrefix string, db *sql.DB) (database.Driver, error) {
	return sqlcipher.WithInstance(db, &sqlcipher.Config{
		MigrationsTable: tablePrefix + sqlcipher.DefaultMigrationsTable,
	})
}
