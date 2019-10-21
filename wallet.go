package chug

import (
	"github.com/crypto-hug/crypto-hug/utils"
	"github.com/v-braun/go-must"
)

type Wallet struct {
	PrivK []byte
	PubK  []byte
	Addr  string
}

func NewWallet() *Wallet {
	priv, pub, err := utils.CreateKeyPair()
	must.NoError(err, "could not gen key pair")

	result := NewWalletFromKeys(priv, pub)
	return result
}

func NewWalletFromKeys(priv, pub []byte) *Wallet {
	var err error
	result := new(Wallet)
	result.PrivK = priv
	result.PubK = pub
	result.Addr, err = NewAddress(pub)
	must.NoError(err, "could not create address")

	return result
}
