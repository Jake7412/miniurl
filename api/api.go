package api

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Handler is defined by the consumer (top-down / caller-owned interface).
type Handler interface {
	AddUrl(url string) (hash string, err error)
}

type API struct {
	handler Handler
}

func Bind(r *httprouter.Router, h Handler) {
	a := &API{handler: h}
	r.POST("/api/v1/url", a.AddUrl)
}

type AddUrlReq struct {
	Url string `json:"url"`
}

type AddUrlResp struct {
	Url  string `json:"url"`
	Hash string `json:"hash"`
}

type ErrorResp struct {
	Msg string `json:"msg"`
}

func (a *API) AddUrl(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var v AddUrlReq
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		respondError(w, http.StatusBadRequest, "bad request")
		return
	}
	hash, err := a.handler.AddUrl(v.Url)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, AddUrlResp{Url: v.Url, Hash: hash})
}

func respondJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		fmt.Println(err)
	}
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, ErrorResp{Msg: msg})
}
