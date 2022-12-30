package sqlite

import (
	"database/sql"

	"github.com/richard-ramos/komainu/pkg/persistence"

	_ "github.com/mutecomm/go-sqlcipher/v4"
)

func Initialize(dburl string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", URLDefaults(dburl))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	driver, err := Driver(db)
	if err != nil {
		return nil, err
	}

	err = persistence.Migrate(AssetNames(), Asset, driver, "sqlite3")
	if err != nil {
		return nil, err
	}

	return db, nil
}
