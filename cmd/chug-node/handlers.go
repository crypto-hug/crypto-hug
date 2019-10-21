package main

import (
	"errors"
	"net/http"

	chug "github.com/crypto-hug/crypto-hug"
	"github.com/crypto-hug/crypto-hug/api"
	"github.com/crypto-hug/crypto-hug/cmd/chug-node/client"
	"github.com/gorilla/mux"
)

func (nh *nodeHost) postTx(res *api.Response, req *api.Request) {
	host := nh.host

	apiModel := &client.Tx{}

	err := req.JSONRequest(apiModel)
	nh.MustNoError(err, http.StatusBadRequest, "invalid body")

	tx, err := mapTx(apiModel)
	nh.MustNoError(err, http.StatusBadRequest, err)

	err = host.ProcessTransaction(tx)
	nh.MustNoError(err, http.StatusBadRequest, err)

	res.EmptyResponse(http.StatusAccepted)
}

func (nh *nodeHost) getHugEtag(res *api.Response, req *api.Request) {
	vars := mux.Vars(req.Request)
	addr, ok := vars["address"]
	if !ok {
		nh.MustNoError(errors.New("missing route parameter 'address'"), http.StatusBadRequest, nil)
	}

	result, err := nh.host.GetHugEtag(addr)
	nh.MustNoError(err, http.StatusBadRequest, nil)

	res.JSONRespnse(http.StatusOK, &struct{ Etag string }{Etag: result})
}

func (nh *nodeHost) getVersion(res *api.Response, req *api.Request) {
	result := &client.Ver{
		Blockchain: string(chug.BlockchainVersion),
	}

	res.JSONRespnse(http.StatusOK, result)
}
