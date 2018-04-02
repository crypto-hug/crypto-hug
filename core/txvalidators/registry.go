package txvalidators

import (
	".."
	"../errors"
)

var me *CommonValidatorRegistry = nil

type CommonValidatorRegistry struct {
	validators map[core.TransactionType]core.TxValidators
}

func init() {
	me = &CommonValidatorRegistry{}
	me.validators = make(map[core.TransactionType]core.TxValidators)
	common := &CommonValidator{}
	spawnHug := &SpawnHugValidator{}

	me.validators[core.SpawnHugTxType] = core.TxValidators{common, spawnHug}
}

func SharedRegistry() *CommonValidatorRegistry {
	return me
}

func (self *CommonValidatorRegistry) Get(tx *core.Transaction) (core.TxValidators, error) {
	list, has := self.validators[tx.Type]
	if !has {
		return nil, errors.TxValidationUnknownTransactionType(string(tx.Type))
	}

	return list, nil
}
