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
)

func (h *Handler) walletBalance(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)
	log.Info("api.wallet.balance.req")

	// UID.
	uid, _, err := h.userinfo("walletBalance", r)
	if err != nil {
		log.Error("api.wallet.balance.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Balance.
	balance, err := wdb.Balance(uid)
	if err != nil {
		log.Error("api.wallet.balance.wdb.balance.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Response.
	rsp := proto.WalletBalanceResponse{
		AllBalance:         balance.AllBalance,
		UnconfirmedBalance: balance.UnconfirmedBalance,
	}
	resp.writeJSON(rsp)
}

func (h *Handler) walletUnspent(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	// UID.
	uid, _, err := h.userinfo("walletUnspent", r)
	if err != nil {
		log.Error("api.wallet.unspent.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.WalletUnspentRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].unspent.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet.unspent.req:%+v", req)

	rsp, err := wdb.Unspents(uid, req.Amount)
	if err != nil {
		log.Error("api.wallet[%v].unspent.by.amount.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet.unspent.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

func (h *Handler) walletPushTx(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	log.Info("api.wallet.push.tx.req:%+v", nil)
	// UID.
	uid, _, err := h.userinfo("walletPushTx", r)
	if err != nil {
		log.Error("api.wallet.push.tx.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.TxPushRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].push.tx.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet.push.tx.req:%+v", req)

	txid, err := wdb.chain.PushTx(req.TxHex)
	if err != nil {
		log.Error("api.wallet[%v].push.tx.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	rsp := &proto.TxPushResponse{
		TxID: txid,
	}
	resp.writeJSON(rsp)
}
