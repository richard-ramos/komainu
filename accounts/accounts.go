package accounts

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"

	"github.com/google/uuid"
)

type Keypair struct {
	Private ed25519.PrivateKey
	Public  ed25519.PublicKey
}

type Account struct {
	id      string
	keypair Keypair
}

func NewAccount(privateKey []byte) (Account, error) {
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

	uuid, err := uuid.NewRandom()
	if err != nil {
		return Account{}, err
	}

	return Account{
		id:      uuid.String(),
		keypair: keypair,
	}, nil
}
