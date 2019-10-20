package chug

import (
	"github.com/crypto-hug/crypto-hug/fs"
	"github.com/pkg/errors"
)

type txProcessCtx struct {
	tx *Transaction

	address          string
	issuerAddress    string
	validatorAddress string
	isGenesisTx      bool
}

type TxProcessor struct {
	config  *Config
	fs      *fs.FileSystem
	states  *StateStore
	txStore *TxStore
}

func NewTxProcessor(fs *fs.FileSystem, config *Config) *TxProcessor {
	proc := new(TxProcessor)
	proc.fs = fs
	proc.config = config
	proc.states = newStateStore(fs, config)
	proc.txStore = newTxStore(fs, config)

	return proc
}

func (proc *TxProcessor) Process(tx *Transaction) error {
	ctx, addr, err := proc.newTxProcessCtx(tx, proc.config)
	if err != nil {
		return errors.Wrapf(err, "tx %s failed validation", addr)
	}

	if err := proc.validate(ctx); err != nil {
		return errors.Wrapf(err, "validation of tx [%s] failed", addr)
	}

	if ctx.isGenesisTx {
		err = proc.processGenesisTx(ctx)
	} else if ctx.tx.Type == GiveHugTransactionType {
		err = proc.processGiveHugTx(ctx)
	} else {
		err = errors.Errorf("unknwon tx type {%s}", ctx.tx.Type)
	}

	if err != nil {
		return err
	}

	if ctx.isGenesisTx || proc.txStore.StagedTxCount() >= proc.config.Blocks.Size {
		proc.txStore.CommitStagedTx()
	}

	return nil
}

func (proc *TxProcessor) processGenesisTx(ctx *txProcessCtx) error {
	if err := proc.states.HugCreateIfNotExists(ctx.issuerAddress); err != nil {
		return err
	}
	if err := proc.states.HugSetEtag(ctx.issuerAddress, ctx.address); err != nil {
		return err
	}
	if err := proc.txStore.StageTx(ctx); err != nil {
		return err
	}

	return nil
}

func (proc *TxProcessor) processGiveHugTx(ctx *txProcessCtx) error {
	if err := proc.states.HugCreateIfNotExists(ctx.validatorAddress); err != nil {
		return err
	}

	if err := proc.states.HugAddLinks(ctx); err != nil {
		return err
	}

	if err := proc.states.HugSetEtag(ctx.issuerAddress, ctx.address); err != nil {
		return err
	}

	if err := proc.states.HugSetEtag(ctx.validatorAddress, ctx.address); err != nil {
		return err
	}

	if err := proc.txStore.StageTx(ctx); err != nil {
		return err
	}

	return nil
}

func (proc *TxProcessor) validate(ctx *txProcessCtx) error {
	if err := ctx.tx.Check(); err != nil {
		return errors.Wrapf(err, "tx [%s] validation failed check", ctx.address)
	}

	if !proc.states.HugExists(ctx.issuerAddress) && !ctx.isGenesisTx {
		return errors.Errorf("issuer hug [%s] does not exists", ctx.issuerAddress)
	}

	if ctx.issuerAddress == ctx.validatorAddress && !ctx.isGenesisTx {
		return errors.Errorf("self hugging is not possible for address [%s]", ctx.issuerAddress)
	}

	issuerEtag, err := proc.states.HugGetEtag(ctx.issuerAddress)
	if err != nil {
		return errors.Wrapf(err, "issuer hug [%s] failed get etag", ctx.issuerAddress)
	}
	if issuerEtag == "" && !ctx.isGenesisTx {
		return errors.Errorf("issuer hug [%s] has empty etag", ctx.issuerAddress)
	}
	if issuerEtag != ctx.tx.IssuerEtag {
		return errors.Errorf("actual issuer etag [%s] does not match with tx issuer etag [%s]", issuerEtag, ctx.tx.IssuerEtag)
	}

	validatorEtag, err := proc.states.HugGetEtag(ctx.validatorAddress)
	if err != nil {
		return errors.Wrapf(err, "validator hug [%s] failed get etag", ctx.validatorAddress)
	}
	if validatorEtag != ctx.tx.ValidatorEtag {
		return errors.Errorf("actual validator etag [%s] does not match with tx validator etag [%s]", validatorEtag, ctx.tx.ValidatorEtag)
	}

	return nil
}

func (proc *TxProcessor) newTxProcessCtx(tx *Transaction, conf *Config) (*txProcessCtx, string, error) {
	var err error
	result := new(txProcessCtx)
	result.address, err = tx.Address()
	if err != nil {
		return nil, "unknown", errors.Wrap(err, "could not create address for tx")
	}

	result.issuerAddress, err = tx.IssuerAddress()
	if err != nil {
		return nil, result.address, errors.Wrap(err, "could not create issuer address")
	}

	result.validatorAddress, err = tx.ValidatorAddress()
	if err != nil {
		return nil, result.address, errors.Wrap(err, "could not create validator address")
	}

	result.tx = tx
	result.isGenesisTx = tx.IsGenesisTx(proc.config)

	return result, result.address, nil
}
