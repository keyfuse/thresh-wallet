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
	"time"

	"proto"

	jwt "github.com/dgrijalva/jwt-go"
)

func (h *Handler) loginVCode(w http.ResponseWriter, r *http.Request) {
	log := h.log
	vcode := h.loginCode
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

func (h *Handler) loginToken(w http.ResponseWriter, r *http.Request) {
	log := h.log
	conf := h.conf
	vcode := h.loginCode
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
	if conf.EnableVCode {
		if err := vcode.Check(req.UID, req.VCode); err != nil {
			log.Error("api.token.vcode.error:%+v", err)
			resp.writeErrorWithStatus(400, err)
			return
		}
		vcode.Remove(req.UID)
	}

	// Make token.
	_, token, err := tokenAuth.Encode(jwt.MapClaims{"uid": req.UID, "did": req.DeviceID, "t": time.Now().Unix(), "net": conf.ChainNet})
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
