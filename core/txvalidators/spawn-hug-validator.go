package txvalidators

import (
	".."
	"../chug"
	"../errors"
)

type SpawnHugValidator struct {
}

func (self SpawnHugValidator) Validate(tx *core.Transaction) error {
	err := errors.MustBeNotNil("tx", tx)
	if err != nil {
		return err
	}

	if tx.Type != chug.SpawnHugTxType {
		return errors.NewErrorFromString("Unknown tx type %s", tx.Type)
	}

	return nil
}
