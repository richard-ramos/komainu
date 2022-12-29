package sqlite

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
)

func Driver(tableName string, db *sql.DB) (database.Driver, error) {
	return sqlite3.WithInstance(db, &sqlite3.Config{
		MigrationsTable: tableName + sqlite3.DefaultMigrationsTable,
	})
}
