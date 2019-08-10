// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package server

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"proto"

	jwt "github.com/dgrijalva/jwt-go"
)

func (h *Handler) loginVCode(w http.ResponseWriter, r *http.Request) {
	log := h.log
	smtp := h.smtp
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

	result, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		resp.writeError(fmt.Errorf("api.generate.vcode.error"))
		return
	}

	code := fmt.Sprintf("%06v", result)
	vcode.Add(req.UID, code)
	switch loginType(req.UID) {
	case Mobile:
		log.Info("api.vcode.mobile.resp:[uid:%v, vcode:%v]", req.UID, code)
	case Email:
		log.Info("api.vcode.email.resp:[uid:%v, vcode:%v]", req.UID, code)
		if err := smtp.VCode(req.UID, "KeyFuse Labs", code); err != nil {
			log.Error("api.vcode.email.send.error:%+v", err)
			resp.writeError(fmt.Errorf("api.email.send.vcode.error"))
			return
		}
	default:
		log.Error("api.vcode.uid.type.unknow:%+v", req.UID)
		resp.writeError(fmt.Errorf("api.vcode.uid.type.unknow:%v, need.mobile.or.email", req.UID))
		return
	}
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
