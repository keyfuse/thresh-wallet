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
	log        *xlog.Log
	wdb        *WalletDB
	conf       *Config
	smtp       *Smtp
	netprefix  string
	tokenAuth  *jwtauth.JWTAuth
	loginCode  *Vcode
	backupCode *Vcode
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
	smtp := NewSmtp(log, conf)
	loginCode := NewVcode(log, conf)
	backupCode := NewVcode(log, conf)
	tokenAuth := jwtauth.New("HS256", []byte(conf.TokenSecret), nil)
	handler := &Handler{
		log:        log,
		wdb:        wdb,
		conf:       conf,
		smtp:       smtp,
		loginCode:  loginCode,
		backupCode: backupCode,
		netprefix:  netprefix,
		tokenAuth:  tokenAuth,
	}
	return handler
}

// Init -- starts the handler.
func (h *Handler) Init() error {
	conf := h.conf
	wdb := h.wdb
	return wdb.Open(conf.DataDir)
}

// Close -- used to close the handler.
func (h *Handler) Close() {
	wdb := h.wdb
	wdb.Close()
}

func (h *Handler) userinfo(tag string, r *http.Request) (string, error) {
	log := h.log

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Error("api.handler[%v].uid.jwtauth.error:%+v", tag, err)
		return "", err
	}
	return fmt.Sprintf("%v", claims["uid"]), nil
}
