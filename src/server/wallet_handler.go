// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
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
