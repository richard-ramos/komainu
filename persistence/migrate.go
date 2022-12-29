package persistence

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

const tbl_prefix = "komainu_"

type assetFn = func(name string) ([]byte, error)
type assets = []string

// Migrate applies migrations.
func Migrate(assetNames assets, assetFn assetFn, driver database.Driver, databaseName string) error {
	return migrateDB(driver, databaseName, bindata.Resource(
		assetNames,
		assetFn,
	))
}

// Migrate database using provided resources.
func migrateDB(driver database.Driver, databaseName string, resources *bindata.AssetSource) error {
	source, err := bindata.WithInstance(resources)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"go-bindata",
		source,
		databaseName,
		driver)
	if err != nil {
		return err
	}

	if err = m.Up(); err != migrate.ErrNoChange {
		return err
	}
	return nil
}
