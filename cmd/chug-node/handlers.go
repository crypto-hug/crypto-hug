package main

import (
	"errors"
	"net/http"

	chug "github.com/crypto-hug/crypto-hug"
	"github.com/crypto-hug/crypto-hug/cmd/chug-node/client"
	"github.com/gorilla/mux"
)

func (api *api) postTx(res *Response, req *Request) {
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

func (api *api) getHugEtag(res *Response, req *Request) {
	vars := mux.Vars(req.Request)
	addr, ok := vars["address"]
	if !ok {
		api.PanicWhenError(errors.New("missing route parameter 'address'"), http.StatusBadRequest, nil)
	}

	result, err := api.host.GetHugEtag(addr)
	api.PanicWhenError(err, http.StatusBadRequest, nil)

	res.JSONRespnse(http.StatusOK, &struct{ Etag string }{Etag: result})
}

func (api *api) getVersion(res *Response, req *Request) {
	result := &client.Ver{
		Blockchain: string(chug.BlockchainVersion),
	}

	res.JSONRespnse(http.StatusOK, result)
}
