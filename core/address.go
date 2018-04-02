package core

import (
	"../crypt"
	"../formatters"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

var versionByte = byte(AddressVersion)

type Address struct {
	Address    string
	PubKeyHash string
}

func NewAddressFromRaw(raw []byte) *Address {
	versionPos := 1
	checksumPos := len(raw) - addrCheckSumLen
	pubKeyHashRaw := raw[versionPos:checksumPos]

	address := formatters.Base58String(raw)
	pubKeyHash := formatters.Base58String(pubKeyHashRaw)

	result := Address{Address: address, PubKeyHash: pubKeyHash}

	return &result
}

func NewAddress() (*Address, error) {
	randomId := crypt.NewId()
	result, err := NewAddressFromPubKey(randomId)
	return result, err
}

func NewAddressFromPubKey(pub []byte) (*Address, error) {
	pubHashRaw, err := createPubKeyHash(pub)
	if err != nil {
		return nil, err
	}

	pubHashVersioned := append([]byte{versionByte}, pubHashRaw...)
	checksum := versionedPubHashChecksum(pubHashVersioned)
	addressRaw := append(pubHashVersioned, checksum...)

	addr := formatters.Base58String(addressRaw)
	pubHash := formatters.Base58String(pubHashRaw)

	result := Address{Address: addr, PubKeyHash: pubHash}

	return &result, err
}

func NewAddressFromString(address string) (*Address, error) {
	raw, err := formatters.Base58FromString(address)
	if err != nil {
		return nil, err
	}

	result := NewAddressFromRaw(raw)

	return result, err
}

func versionedPubHashChecksum(verPubKey []byte) []byte {
	result := sha256.Sum256(verPubKey)
	result = sha256.Sum256(result[:])

	return result[:addrCheckSumLen]
}

func createPubKeyHash(pubKey []byte) ([]byte, error) {
	sha256Checksum := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(sha256Checksum[:])
	if err != nil {
		return nil, err
	}

	result := RIPEMD160Hasher.Sum(nil)

	return result, nil
}
