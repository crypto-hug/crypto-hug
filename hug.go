package chug

import (
	"crypto/sha256"

	"github.com/crypto-hug/crypto-hug/utils"
	"golang.org/x/crypto/ripemd160"
)

const hugAddrVersion = 1
const addrCheckSumLen = 4
const versionByte = byte(hugAddrVersion)

func HugAddrFromPubKey(pub []byte) (string, error) {
	pubHashRaw, err := createPubKeyHash(pub)
	if err != nil {
		return "", err
	}

	pubHashVersioned := append([]byte{versionByte}, pubHashRaw...)
	checksum := versionedPubHashChecksum(pubHashVersioned)
	addressRaw := append(pubHashVersioned, checksum...)

	addr := utils.Base58ToStr(addressRaw)

	return addr, err
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
