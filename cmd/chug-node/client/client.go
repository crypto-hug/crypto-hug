package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	chug "github.com/crypto-hug/crypto-hug"
	"github.com/pkg/errors"

	"github.com/v-braun/go-must"
)

type Client struct {
	*http.Client
}

func NewAPIClient() *Client {
	c := &Client{}

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/tmp/chug.sock")
			},
		},
	}

	c.Client = &httpc

	return c
}

func (c *Client) PostJson(path string, data interface{}, result interface{}) error {
	body, err := json.Marshal(data)
	must.NoError(err, "unexpected error during marshalling")

	r, err := c.Post("http://unix"+path, "application/json", bytes.NewBuffer(body))
	err = c.handleResponse(path, r, err)
	if err != nil {
		return err
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return nil
	}

	rbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.WithStack(err)
	}

	err = json.Unmarshal(rbody, result)
	return errors.WithStack(err)
}

func (c *Client) handleResponse(path string, r *http.Response, err error) error {
	if err != nil {
		return errors.WithStack(err)
	}

	if r.StatusCode >= 300 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return errors.WithStack(err)
		}

		msg := r.Status
		if len(body) > 0 {
			msg = string(body)
		}

		return errors.New(fmt.Sprintf("request to %s failed | code [%d] | message [%s] ", path, r.StatusCode, msg))
	}

	return nil
}

func (c *Client) GetJson(path string, result interface{}) error {
	r, err := c.Get("http://unix" + path)
	err = c.handleResponse(path, r, err)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, result)
	return err
}

func (cli *Client) MustGetVersion() *Ver {
	res := &Ver{}
	err := cli.GetJson("/version", res)

	must.NoError(err, "could not get version")
	return res
}

func (cli *Client) PostTransaction(tx *chug.Transaction) error {
	err := cli.PostJson("/tx", tx, nil)
	return err
}

func (cli *Client) GetHugEtag(addr string) (string, error) {
	res := struct{ Etag string }{Etag: ""}
	err := cli.GetJson("/hug/"+addr+"/etag", &res)

	return res.Etag, err
}
