// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"proto"

	jwt "github.com/dgrijalva/jwt-go"
)

func (h *Handler) vcodefn(w http.ResponseWriter, r *http.Request) {
	log := h.log
	vcode := h.vcode
	resp := newResponse(log, w)

	// Request.
	req := &proto.VCodeRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.vcode.decode.body.error:%+v", err)
		resp.writeError(err)
		return
	}
	log.Info("api.vcode.req:%+v", req)

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	vcode.Add(req.UID, code)
	log.Info("api.vcode.resp:[uid:%v, vcode:%v]", req.UID, code)
}

func (h *Handler) tokenfn(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	conf := h.conf
	vcode := h.vcode
	resp := newResponse(log, w)
	tokenAuth := h.tokenAuth

	// Request.
	req := &proto.TokenRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.token.decode.body.error:%+v", err)
		resp.writeError(err)
		return
	}
	log.Info("api.token.req:%+v", req)

	// Check Vcode.
	if !conf.DisableVCode {
		if err := vcode.Check(req.UID, req.VCode); err != nil {
			log.Error("api.token.vcode.error:%+v", err)
			resp.writeError(err)
			return
		}
		vcode.Remove(req.UID)
	}

	// Check chainnet.
	if !strings.HasPrefix(req.MasterPubKey, h.netprefix) {
		log.Error("api.token.req.masterpubkey[%v].chainnet.error.server:%s", req.MasterPubKey, h.netprefix)
		resp.writeError(fmt.Errorf("chainnet.error.server.is:%v", h.netprefix))
		return
	}

	// Check wallet.
	_, err = wdb.OpenUIDWallet(req.UID, req.MasterPubKey)
	if err != nil {
		log.Error("api.token.open.wallet.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Make token.
	_, token, err := tokenAuth.Encode(jwt.MapClaims{"uid": req.UID, "mpk": req.MasterPubKey, "did": req.DeviceID, "t": time.Now().Unix(), "net": conf.ChainNet})
	if err != nil {
		log.Error("api.token[%+v].error:%+v", req, err)
		resp.writeError(err)
		return
	}

	// Response.
	rsp := proto.TokenResponse{
		Token: token,
	}
	resp.writeJSON(rsp)
}
