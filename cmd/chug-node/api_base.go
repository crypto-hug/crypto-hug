package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiError struct {
	Code    int
	Message string
}

func (err *ApiError) Error() string { return err.Message }

type ApiBase struct {
	*mux.Router
	Listener net.Listener
}
type Response struct {
	http.ResponseWriter
}
type Request struct {
	*http.Request
}

func (a *ApiBase) Post(path string, f func(w *Response, r *Request)) {
	handler := a.wrapHandleFunc(f)
	a.Router.HandleFunc(path, handler).Methods("POST")
}

func (a *ApiBase) Put(path string, f func(w *Response, r *Request)) {
	handler := a.wrapHandleFunc(f)
	a.Router.HandleFunc(path, handler).Methods("PUT")
}
func (a *ApiBase) Get(path string, f func(w *Response, r *Request)) {
	handler := a.wrapHandleFunc(f)
	a.Router.HandleFunc(path, handler).Methods("GET")
}

func (a *ApiBase) PanicWhenError(err error, code int, msg interface{}) {
	if err == nil {
		return
	}

	panic(&ApiError{
		Code:    code,
		Message: fmt.Sprintf("%s", msg),
	})
}

func (a *ApiBase) wrapHandleFunc(handler func(w *Response, r *Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res := &Response{
			ResponseWriter: w,
		}
		req := &Request{
			Request: r,
		}

		defer func() {
			if err := recover(); err != nil {
				if apiErr, ok := err.(*ApiError); ok {
					res.ErrorResponse(apiErr.Code, apiErr.Error())
				} else {
					msg := "unexpected server error"
					if commonErr, ok := err.(error); ok {
						msg = commonErr.Error()
					}

					res.ErrorResponse(http.StatusInternalServerError, msg)
				}
			}
		}()

		handler(res, req)
	}
}

func (r *Request) JsonBody(out interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, out); err != nil {
		return err
	}

	return nil
}

func (w *Response) JSONRespnse(status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func (w *Response) EmptyResponse(status int) {
	w.WriteHeader(status)
}

func (w *Response) ErrorResponse(code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(message))
	// w.JSONRespnse(code, map[string]string{"error": message})
}
