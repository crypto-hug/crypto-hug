package serialization

import (
	"bytes"
	"encoding/gob"
)

func SimpleEncode(data interface{}) (*bytes.Buffer, error) {
	var result bytes.Buffer
	var encoder = gob.NewEncoder(&result)
	var err = encoder.Encode(data)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
