package crypt

import (
	uuid "github.com/google/uuid"
)

func NewId() []byte{
	id, _ := uuid.New().MarshalBinary()
	return id
}