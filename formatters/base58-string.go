package formatters

import (
	"github.com/mr-tron/base58/base58"
	"github.com/pkg/errors"
)

func Base58String(data []byte) string {
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
