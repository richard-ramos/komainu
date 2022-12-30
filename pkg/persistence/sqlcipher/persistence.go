package sqlcipher

import (
	"database/sql"

	"github.com/richard-ramos/komainu/pkg/persistence"

	_ "github.com/mutecomm/go-sqlcipher/v4"
)

const DefaultKDFIterations = 256_000

func Initialize(dburl string, key string, kdfIterationsNumber int) (*sql.DB, error) {
	if kdfIterationsNumber <= 0 {
		kdfIterationsNumber = DefaultKDFIterations
	}

	db, err := sql.Open("sqlite3", URLDefaults(dburl, key, kdfIterationsNumber))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	driver, err := Driver(db)
	if err != nil {
		return nil, err
	}

	err = persistence.Migrate(AssetNames(), Asset, driver, "sqlcipher")
	if err != nil {
		return nil, err
	}

	return db, nil
}
