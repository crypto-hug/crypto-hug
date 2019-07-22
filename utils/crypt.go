package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
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

func CreateKeyPair() (private []byte, public []byte, err error) {
	curve := defaultCurve()

	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	priv := privKey.D.Bytes()

	x := privKey.PublicKey.X.Bytes()
	x = append(bytes.Repeat([]byte{0x00}, 32-len(x)), x...)

	y := privKey.PublicKey.Y.Bytes()
	y = append(bytes.Repeat([]byte{0x00}, 32-len(y)), y...)

	pub := append(x, y...)

	return priv, pub, err
}

func SignCreate(priv []byte, pub []byte, hash []byte) (sig []byte, err error) {
	privKey := unwrapPrivKey(priv, pub)
	r, s, err := ecdsa.Sign(rand.Reader, &privKey, hash)
	if err != nil {
		return nil, err
	}

	rb := r.Bytes()
	rb = append(bytes.Repeat([]byte{0x00}, 32-len(rb)), rb...)

	sb := s.Bytes()
	sb = append(bytes.Repeat([]byte{0x00}, 32-len(sb)), sb...)

	result := append(rb, sb...)

	return result, nil
}

func SignCheck(pub []byte, hash []byte, sig []byte) bool {
	l := len(sig) / 2
	r := big.Int{}
	s := big.Int{}
	r.SetBytes(sig[:l])
	s.SetBytes(sig[l:])

	pubKey := unwrapPubKey(pub)

	result := ecdsa.Verify(&pubKey, hash, &r, &s)
	return result
}

func unwrapPrivKey(priv []byte, pub []byte) ecdsa.PrivateKey {
	pubKey := unwrapPubKey(pub)
	d := big.Int{}
	d.SetBytes(priv)

	privKey := ecdsa.PrivateKey{PublicKey: pubKey, D: &d}

	return privKey
}

func defaultCurve() elliptic.Curve {
	return elliptic.P256()
}

func unwrapPubKey(pub []byte) ecdsa.PublicKey {
	curve := defaultCurve()
	l := len(pub) / 2
	x := big.Int{}
	y := big.Int{}
	x.SetBytes(pub[:l])
	y.SetBytes(pub[l:])

	result := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

	return result
}
