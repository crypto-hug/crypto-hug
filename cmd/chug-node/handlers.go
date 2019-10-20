package main

import (
	"net/http"

	chug "github.com/crypto-hug/crypto-hug"
	"github.com/crypto-hug/crypto-hug/cmd/chug-node/client"
)

func (api *api) txPost(res *Response, req *Request) {
	host := api.host

	apiModel := &client.Tx{}

	err := req.JsonBody(apiModel)
	api.PanicWhenError(err, http.StatusBadRequest, "invalid body")

	tx, err := mapTx(apiModel)
	api.PanicWhenError(err, http.StatusBadRequest, err)

	err = host.ProcessTransaction(tx)
	api.PanicWhenError(err, http.StatusBadRequest, err)

	res.EmptyResponse(http.StatusAccepted)
}

func (api *api) versionGet(res *Response, req *Request) {
	result := &client.Ver{
		Blockchain: string(chug.BlockchainVersion),
	}

	res.JSONRespnse(http.StatusOK, result)
}
