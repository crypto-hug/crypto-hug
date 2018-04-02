package crypt

import (
	"math/big"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

func CreateKeyPair() (private []byte, public []byte, err error){
	curve := defaultCurve()
	
	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil{
		return nil, nil, err
	}

	priv := privKey.D.Bytes()
	pub := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
		
	return priv, pub, err
}
func defaultCurve() elliptic.Curve{
	return elliptic.P256()
}
func unwrapPubKey(pub []byte) ecdsa.PublicKey{
	curve := defaultCurve()
	l := len(pub) / 2
	x := big.Int{}
	y := big.Int{}
	x.SetBytes(pub[:l])
	y.SetBytes(pub[l:])

	result := ecdsa.PublicKey{Curve: curve, X:&x, Y:&y}

	return result
}
func unwrapPrivKey(priv []byte, pub []byte) ecdsa.PrivateKey{
	pubKey := unwrapPubKey(pub)
	d := big.Int{}
	d.SetBytes(priv)

	privKey := ecdsa.PrivateKey{PublicKey: pubKey, D: &d}

	return privKey
}

func SignHash(priv []byte, pub []byte, hash []byte) (sig []byte, err error){
	privKey := unwrapPrivKey(priv, pub)
	r, s, err := ecdsa.Sign(rand.Reader, &privKey, hash)
	if (err != nil){
		return nil, err
	}

	result := append(r.Bytes(), s.Bytes()...)
	return result, nil
}

func VerifySig(pub []byte, hash []byte, sig []byte) bool{
	l := len(sig) / 2
	r := big.Int{}
	s := big.Int{}
	r.SetBytes(sig[:l])
	s.SetBytes(sig[l:])

	pubKey := unwrapPubKey(pub)

	result := ecdsa.Verify(&pubKey, hash, &r, &s)
	return result
}

