package core

import (
	"../crypt"
	"../formatters"
)

type Wallet struct{
	PrivKey []byte
	PubKey []byte
	Address *Address
}


const addrCheckSumLen = 4

func NewWallet() (*Wallet, error){
	priv, pub, err := crypt.CreateKeyPair()
	if err != nil {
		return nil, err
	}

	result := Wallet{PrivKey: priv, PubKey: pub}

	address, err := NewAddressFromPubKey(pub)
	if err != nil {
		return nil, err
	}

	result.Address = address


	return &result, nil
}




func (self *Wallet) PrivAsString() string{
	result := formatters.Base58String(self.PrivKey)
	return result
}

func (self *Wallet) PubAsString() string{
	result := formatters.Base58String(self.PubKey)
	return result
}

