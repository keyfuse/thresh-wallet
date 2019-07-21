// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"encoding/json"
	"net/http"
	"xlog"
)

type response struct {
	log        *xlog.Log
	w          http.ResponseWriter
	StatusCode int
}

func newResponse(log *xlog.Log, w http.ResponseWriter) *response {
	return &response{
		log:        log,
		w:          w,
		StatusCode: http.StatusOK,
	}
}

func (r *response) writeError(err error) {
	w := r.w
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func (r *response) writeErrorWithStatus(status int, err error) {
	w := r.w
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}

func (r *response) writeJSON(thing interface{}) {
	w := r.w
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	if err := encoder.Encode(thing); err != nil {
		r.log.Error("api.write.json.error:%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
