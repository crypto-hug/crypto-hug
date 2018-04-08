package core

type Version int32

const defaultVersion Version = 1

const (
	AssetVersion   Version = defaultVersion
	TxVersion      Version = defaultVersion
	BlockVersion   Version = defaultVersion
	AddressVersion Version = defaultVersion
)
