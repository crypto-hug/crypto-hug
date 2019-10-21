package main

import (
	"net"
	"net/http"
	"os"

	chug "github.com/crypto-hug/crypto-hug"
	"github.com/crypto-hug/crypto-hug/api"
	"github.com/crypto-hug/crypto-hug/cmd/chug-node/client"
	"github.com/crypto-hug/crypto-hug/utils"
	"github.com/pkg/errors"
)

type nodeHost struct {
	*api.Api
	host *chug.NodeHost
}

const sockAddr = "/tmp/chug.sock"

func mapTx(apiModel *client.Tx) (*chug.Transaction, error) {
	tx := &chug.Transaction{}
	tx.Version = chug.Version(apiModel.Version)
	tx.Type = chug.TransactionType(apiModel.Type)
	tx.Timestamp = apiModel.Timestamp
	tx.IssuerEtag = apiModel.IssuerEtag
	tx.ValidatorEtag = apiModel.ValidatorEtag

	var err error
	if tx.Hash, err = utils.NewBase58JsonValFromString(apiModel.Hash); err != nil {
		return nil, errors.Wrap(err, "invalid hash")
	}
	if tx.IssuerPubKey, err = utils.NewBase58JsonValFromString(apiModel.IssuerPubKey); err != nil {
		return nil, errors.Wrap(err, "invalid issuerPubKey")
	}
	if tx.IssuerLock, err = utils.NewBase58JsonValFromString(apiModel.IssuerLock); err != nil {
		return nil, errors.Wrap(err, "invalid issuerLock")
	}

	if tx.ValidatorPubKey, err = utils.NewBase58JsonValFromString(apiModel.ValidatorPubKey); err != nil {
		return nil, errors.Wrap(err, "invalid validatorPubKey")
	}
	if tx.ValidatorLock, err = utils.NewBase58JsonValFromString(apiModel.ValidatorLock); err != nil {
		return nil, errors.Wrap(err, "invalid validatorLock")
	}

	if tx.Data, err = utils.NewBase58JsonValFromString(apiModel.Data); err != nil {
		return nil, errors.Wrap(err, "invalid data")
	}

	return tx, nil
}

func newNodeHost(host *chug.NodeHost) *nodeHost {
	result := new(nodeHost)
	result.Api = api.New()
	result.host = host

	result.Post("/tx", result.postTx)
	result.Get("/version", result.getVersion)
	result.Get("/hug/{address}/etag", result.getHugEtag)

	return result
}

func (nh *nodeHost) Run() error {
	nh.host.Start()
	if err := os.RemoveAll(sockAddr); err != nil {
		panic(errors.Wrapf(err, "failed remove socket %s", sockAddr))
	}

	l, err := net.Listen("unix", sockAddr)
	if err != nil {
		panic(errors.Wrapf(err, "failed listen to socket %s", sockAddr))
	}

	return http.Serve(l, nh)
}
