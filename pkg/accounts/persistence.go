package accounts

import (
	"database/sql"

	"github.com/richard-ramos/komainu/pkg/persistence/sqlcipher"
)

type Persistence struct {
	dbStore *sql.DB
}

func NewPersistence(db *sql.DB) *Persistence {
	return &Persistence{
		dbStore: db,
	}
}

func (p *Persistence) SaveAccount(account Account) error {
	if account.KDFIterations <= 0 {
		account.KDFIterations = sqlcipher.DefaultKDFIterations
	}

	_, err := p.dbStore.Exec("INSERT OR REPLACE INTO accounts (id, kdfIterations) VALUES (?, ?)", account.ID, account.KDFIterations)
	if err != nil {
		return err
	}

	return nil
}
