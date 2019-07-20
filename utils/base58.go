package utils

import (
	"encoding/json"

	"github.com/mr-tron/base58/base58"
	"github.com/pkg/errors"
)

type Base58JsonVal struct {
	val []byte
}

func NewBase58JsonValFromData(val []byte) *Base58JsonVal {
	result := new(Base58JsonVal)
	result.val = val
	return result
}

func NewBase58JsonValFromString(str string) (*Base58JsonVal, error) {
	result := new(Base58JsonVal)
	val, err := Base58FromString(str)
	result.val = val
	return result, err
}

func (s *Base58JsonVal) String() string {
	if s == nil {
		return ""
	}
	return Base58ToStr(s.val)
}
func (s *Base58JsonVal) Bytes() []byte {
	if s == nil {
		return []byte{}
	}
	return s.val
}

func (s *Base58JsonVal) MarshalJSON() ([]byte, error) {
	if s == nil {
		return json.Marshal(nil)
	}
	if len(s.val) <= 0 {
		return json.Marshal(nil)
	}

	data := Base58ToStr(s.val)
	return json.Marshal(data)
}

func (s *Base58JsonVal) UnmarshalJSON(data []byte) error {
	var unmarshalled string
	if err := json.Unmarshal(data, &unmarshalled); err != nil {
		return err
	}

	s.val = []byte{}
	if unmarshalled == "" {
		return nil
	}
	if len(unmarshalled) == 0 {
		return nil
	}

	str, err := Base58FromString(unmarshalled)
	if err != nil {
		return err
	}

	s.val = str
	return nil
}

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
