package utils

import (
	"github.com/mr-tron/base58/base58"
	"github.com/pkg/errors"
)

func Base58ToStr(data []byte) string {
	result := base58.Encode(data)
	return result
}

func Base58FromString(str string) ([]byte, error) {
	result, err := base58.Decode(str)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, err
}

func Base58FromStringMust(str string) []byte {
	result, err := Base58FromString(str)
	if err != nil {
		panic(err)
	}

	return result
}
