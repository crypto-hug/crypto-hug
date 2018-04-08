package core

type AssetType string

type AssetHeader struct {
	Version  Version
	Address  *Address
	Producer *Address
	Type     AssetType
}

func NewAssetHeader(producer *Address, assetType AssetType) *AssetHeader {
	result := AssetHeader{}
	result.Version = AssetVersion
	result.Type = assetType
	result.Address = NewAddressStrict()
	result.Producer = producer

	return &result
}
