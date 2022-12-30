package accounts

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"

	"github.com/google/uuid"
	"github.com/richard-ramos/komainu/pkg/persistence/sqlcipher"
)

type Keypair struct {
	Private ed25519.PrivateKey
	Public  ed25519.PublicKey
}

type Account struct {
	ID            string
	KDFIterations int
	Keypair       Keypair
}

func NewAccount(id string, privateKey []byte) (Account, error) {
	var err error
	var keypair Keypair

	if privateKey == nil {
		keypair.Public, keypair.Private, err = ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return Account{}, err
		}
	} else {
		reader := bytes.NewReader(privateKey)
		keypair.Public, keypair.Private, err = ed25519.GenerateKey(reader)
		if err != nil {
			return Account{}, err
		}
	}

	if id == "" {
		uuid, err := uuid.NewRandom()
		if err != nil {
			return Account{}, err
		}
		id = uuid.String()
	}

	return Account{
		ID:            id,
		Keypair:       keypair,
		KDFIterations: sqlcipher.DefaultKDFIterations,
	}, nil
}
