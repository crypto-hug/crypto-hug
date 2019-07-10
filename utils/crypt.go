package utils

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"

	"github.com/pkg/errors"
)

func Hash(data []byte) []byte {
	var result = sha256.Sum256(data)
	return result[:]
}

func HashAll(allData ...[]byte) []byte {
	var data = bytes.Join(allData, []byte{})
	var result = Hash(data)
	return result
}

func SignCheck(pubKey []byte, hash []byte, signature []byte) error {
	pk, err := pubKeyFromBytes(pubKey)
	if err != nil {
		return err
	}

	err = rsa.VerifyPKCS1v15(pk, crypto.SHA256, hash, signature)
	return err
}

func SignCreate(privKey []byte, hash []byte) ([]byte, error) {
	key, err := privFromBytes(privKey)
	if err != nil {
		return nil, err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash)
	return signature, err
}

func PrivKeyFromString(key string) (*rsa.PrivateKey, error) {
	data, err := Base58FromString(key)
	if err != nil {
		return nil, err
	}

	result, err := privFromBytes(data)
	return result, err
}

func PrivKeyToBytes(pk *rsa.PrivateKey) []byte {
	result := x509.MarshalPKCS1PrivateKey(pk)
	return result
}

func PubKeyToBytes(pk *rsa.PublicKey) []byte {
	result := x509.MarshalPKCS1PublicKey(pk)
	return result
}

func GeneratePrivKey() *rsa.PrivateKey {
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(errors.WithStack(err))
	}

	return k
}

func pubKeyFromBytes(data []byte) (*rsa.PublicKey, error) {
	pub, err := x509.ParsePKCS1PublicKey(data)
	if err != nil {
		return nil, err
	}

	return pub, nil
}

func privFromBytes(data []byte) (*rsa.PrivateKey, error) {
	priv, err := x509.ParsePKCS1PrivateKey(data)
	if err != nil {
		return nil, err
	}

	return priv, nil
}
