package utils

import "encoding/binary"

func Int64GetBytes(val int64) []byte {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, uint64(val))
	return result
}
