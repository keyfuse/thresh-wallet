// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"encoding/json"
	"net/http"

	"proto"
	//"xlog"
	//"github.com/tokublock/tokucore/xcrypto"
)

func (h *Handler) ecdsaNewAddress(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	// UID.
	uid, cliMasterPubKey, err := h.userinfo("ecdsaNewAddress", r)
	if err != nil {
		log.Error("api.ecdsa.newaddress.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.EcdsaAddressRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.ecdsa.newaddress[%v].decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.ecdsa.req:%+v", req)

	// New address.
	address, err := wdb.NewAddressByUID(uid, cliMasterPubKey)
	if err != nil {
		log.Error("api.ecdsa.newaddress.wdb.newaddress.error:%+v", err)
		resp.writeError(err)
		return
	}
	rsp := &proto.EcdsaAddressResponse{
		Pos:     address.Pos,
		Address: address.Address,
	}
	log.Info("api.ecdsa.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

/*
func (h *Handler) ecdsaR2(w http.ResponseWriter, r *http.Request) {
	log := h.log
	resp := newResponse(log, w)

	_, claims, _ := jwtauth.FromContext(r.Context())
	uid := fmt.Sprintf("%v", claims["uid"])
	key, err := h.GetKey(uid)
	if err != nil {
		log.Error("api.ecdsa.r2[%v].error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	req := &proto.EcdsaR2Request{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.ecdsa.r2[%v].decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	hdkey, err := key.Derive(req.IDX)
	if err != nil {
		log.Error("api.ecdsa.r2[%v].hdkey.error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	bobParty := xcrypto.NewEcdsaParty(hdkey.PrivateKey())
	// Skip phase1.
	// Phase2.
	_, _, scalarR2 := bobParty.Phase2(req.Hash)
	shareR := bobParty.Phase3(req.R1)
	rsp := &proto.EcdsaR2Response{
		R2:     scalarR2,
		IDX:    req.IDX,
		ShareR: shareR,
	}
	log.Info("api.ecdsa.r2.rsp:%+v", xlog.Pretty(rsp))
	resp.writeJSON(rsp)
}

func (h *Handler) ecdsaS2(w http.ResponseWriter, r *http.Request) {
	log := h.log
	resp := newResponse(log, w)

	_, claims, _ := jwtauth.FromContext(r.Context())
	uid := fmt.Sprintf("%v", claims["uid"])
	key, err := h.GetKey(uid)
	if err != nil {
		log.Error("api.ecdsa.s2[%v].error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	req := &proto.EcdsaS2Request{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.ecdsa.s2[%v].decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	hdkey, err := key.Derive(req.IDX)
	if err != nil {
		log.Error("api.ecdsa.s2[%v].hdkey.derive.error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	bobParty := xcrypto.NewEcdsaParty(hdkey.PrivateKey())
	// Skip phase1.
	// Phase2.
	bobParty.Phase2(req.Hash)
	shareR := bobParty.Phase3(req.R1)
	if shareR.X.Cmp(req.ShareR.X) != 0 || shareR.Y.Cmp(req.ShareR.Y) != 0 {
		log.Error("api.ecdsa.s2[%v].share.not.equal[%+v].[%+v]", uid, xlog.Pretty(shareR), xlog.Pretty(req.ShareR))
		resp.writeError(fmt.Errorf("api.ecdsa.s2.shareR.not.equal"))
		return
	}
	// Skip phase3.

	// Phase4.
	sig2, err := bobParty.Phase4(req.EncPK1, req.EncPub1, req.ShareR)
	rsp := &proto.EcdsaS2Response{
		S2:  sig2,
		IDX: req.IDX,
	}
	log.Info("api.ecdsa.s2.rsp:%+v", xlog.Pretty(rsp))
	resp.writeJSON(rsp)
}
*/
