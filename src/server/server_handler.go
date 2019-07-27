// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"net/http"
	"time"

	"proto"
)

func (h *Handler) serverInfo(w http.ResponseWriter, r *http.Request) {
	conf := h.conf
	log := h.log

	resp := newResponse(log, w)
	rsp := &proto.ServerInfoResponse{
		ChainNet:    conf.ChainNet,
		ServerTime:  time.Now().UTC().Unix(),
		EnableVCode: conf.EnableVCode,
	}
	resp.writeJSON(rsp)
}
