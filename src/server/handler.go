// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"fmt"
	"net/http"

	"xlog"

	"github.com/go-chi/jwtauth"
)

// Handler --
type Handler struct {
	log       *xlog.Log
	conf      *Config
	wdb       *WalletDB
	vcode     *Vcode
	netprefix string
	tokenAuth *jwtauth.JWTAuth
}

// NewHandler -- creates new Handler.
func NewHandler(log *xlog.Log, conf *Config) *Handler {
	var netprefix string
	switch conf.ChainNet {
	case testnet:
		netprefix = "tpub"
	case mainnet:
		netprefix = "xpub"
	}
	wdb := NewWalletDB(log, conf)
	if err := wdb.Open(conf.DataDir); err != nil {
		log.Panic("handler.new.panic:%+v", err)
	}
	vcode := NewVcode(log, conf)
	tokenAuth := jwtauth.New("HS256", []byte(conf.TokenSecret), nil)
	handler := &Handler{
		log:       log,
		conf:      conf,
		wdb:       wdb,
		vcode:     vcode,
		netprefix: netprefix,
		tokenAuth: tokenAuth,
	}
	return handler
}

func (h *Handler) userinfo(tag string, r *http.Request) (string, string, error) {
	log := h.log

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Error("api.handler[%v].uid.jwtauth.error:%+v", tag, err)
		return "", "", err
	}
	return fmt.Sprintf("%v", claims["uid"]), fmt.Sprintf("%v", claims["mpk"]), nil
}
