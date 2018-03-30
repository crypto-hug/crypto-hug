package crypt

import (
	"bytes"
	"crypto/sha256"
)

func BytesHash(data []byte) []byte {
	var result = sha256.Sum256(data)
	return result[:]
}


func AllBytesHash(allData ...[]byte) []byte {
	var data = bytes.Join(allData, []byte{})
	var result = BytesHash(data)
	return result
}
