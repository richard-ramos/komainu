package persistence

import (
	"database/sql"
	"time"

	"go.uber.org/zap"
)

// DBStore is a MessageProvider that has a *sql.DB connection
type DBStore struct {
	db               *sql.DB
	log              *zap.Logger
	enableMigrations bool
}

// DBOption is an optional setting that can be used to configure the DBStore
type DBOption func(*DBStore) error

// WithDB is a DBOption that lets you use any custom *sql.DB with a DBStore.
func WithDB(db *sql.DB) DBOption {
	return func(d *DBStore) error {
		d.db = db
		return nil
	}
}

type ConnectionPoolOptions struct {
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
}

// WithDriver is a DBOption that will open a *sql.DB connection
func WithDriver(driverName string, datasourceName string, connectionPoolOptions ...ConnectionPoolOptions) DBOption {
	return func(d *DBStore) error {
		db, err := sql.Open(driverName, datasourceName)
		if err != nil {
			return err
		}

		if len(connectionPoolOptions) != 0 {
			db.SetConnMaxIdleTime(connectionPoolOptions[0].ConnectionMaxIdleTime)
			db.SetConnMaxLifetime(connectionPoolOptions[0].ConnectionMaxLifetime)
			db.SetMaxIdleConns(connectionPoolOptions[0].MaxIdleConnections)
			db.SetMaxOpenConns(connectionPoolOptions[0].MaxOpenConnections)
		}

		d.db = db
		return nil
	}
}

// WithMigrationsEnabled is a DBOption used to determine whether migrations should
// be executed or not
func WithMigrationsEnabled(enabled bool) DBOption {
	return func(d *DBStore) error {
		d.enableMigrations = enabled
		return nil
	}
}

func DefaultOptions() []DBOption {
	return []DBOption{
		WithMigrationsEnabled(true),
	}
}
