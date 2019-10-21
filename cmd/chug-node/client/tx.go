package client

type Tx struct {
	Version string `json:"version,omitempty"`

	Type string `json:"type,omitempty"`

	Timestamp int64 `json:"timestamp,omitempty"`

	Hash string `json:"hash,omitempty"`

	IssuerPubKey string `json:"issuerPubKey,omitempty"`

	IssuerLock string `json:"issuerLock,omitempty"`

	IssuerEtag string `json:"issuerEtag,omitempty"`

	ValidatorPubKey string `json:"validatorPubKey,omitempty"`

	ValidatorLock string `json:"validatorLock,omitempty"`

	ValidatorEtag string `json:"validatorEtag,omitempty"`

	Data string `json:"data,omitempty"`
}
