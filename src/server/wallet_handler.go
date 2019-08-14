// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package server

import (
	"encoding/json"
	"net/http"

	"proto"
)

func (h *Handler) walletNewAddress(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("walletNewAddress", r)
	if err != nil {
		log.Error("api.wallet.newaddress.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.WalletNewAddressRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet.newaddress[%v].decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet.newaddress.req:%+v", req)

	// New address.
	address, err := wdb.NewAddress(uid, req.Type)
	if err != nil {
		log.Error("api.wallet.newaddress.wdb.newaddress.error:%+v", err)
		resp.writeError(err)
		return
	}
	rsp := &proto.WalletNewAddressResponse{
		Pos:     address.Pos,
		Address: address.Address,
	}
	log.Info("api.wallet.newaddress.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

func (h *Handler) walletCheck(w http.ResponseWriter, r *http.Request) {
	var walletExists bool
	var backupExists bool
	var backupTimestamp int64
	var backupCloudService string

	log := h.log
	wdb := h.wdb
	conf := h.conf
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("walletCheck", r)
	if err != nil {
		log.Error("api.wallet.check.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.WalletCheckRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].check.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet[%v].check.req:%+v", uid, req)

	wallet := wdb.Wallet(uid)
	if wallet != nil {
		walletExists = true
		backupExists = (wallet.Backup.EncryptedPrvKey != "")
		backupTimestamp = wallet.Backup.Time
		backupCloudService = wallet.Backup.CloudService
	}
	// Response.
	rsp := proto.WalletCheckResponse{
		WalletExists:       walletExists,
		BackupExists:       backupExists,
		ForceRecover:       conf.ForceRecover,
		BackupTimestamp:    backupTimestamp,
		BackupCloudService: backupCloudService,
	}
	log.Info("api.wallet.check.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

func (h *Handler) walletCreate(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	smtp := h.smtp
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("walletCreate", r)
	if err != nil {
		log.Error("api.wallet[%v].create.uid.error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.WalletCreateRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].create.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet.create.req:%+v", req)

	// Verify the pub/prv key pairing.
	if err := verifyPubKey(req.MasterPubKey, req.Signature); err != nil {
		log.Error("api.wallet[%v].create.verify.pubkey.error:%+v", uid, err)
		resp.writeErrorWithStatus(400, err)
		return
	}

	// Create wallet.
	if err := wdb.CreateWallet(uid, req.MasterPubKey); err != nil {
		log.Error("api.wallet[%v].wdb.create.error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	// smtp backup.
	if err := smtp.Backup(uid, "KeyFuse Labs-Server-Wallet-Create"); err != nil {
		log.Error("api.wallet[%v].create.smtp.backup.error:%+v", uid, err)
		resp.writeErrorWithStatus(500, nil)
		return
	}

	// Response.
	rsp := proto.WalletCreateResponse{}
	log.Info("api.wallet.create.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

func (h *Handler) walletBalance(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("walletBalance", r)
	if err != nil {
		log.Error("api.wallet.balance.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.WalletBalanceRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].balance.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet[%v].balance.req:%+v", uid, req)

	// Balance.
	balance, err := wdb.Balance(uid)
	if err != nil {
		log.Error("api.wallet.balance.wdb.balance.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Response.
	rsp := proto.WalletBalanceResponse{
		CoinValue: balance.TotalBalance,
	}
	resp.writeJSON(rsp)
}

func (h *Handler) walletUnspent(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("walletUnspent", r)
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

	unspents, err := wdb.Unspents(uid, req.Amount)
	if err != nil {
		log.Error("api.wallet[%v].unspent.by.amount.error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	var rsp []proto.WalletUnspentResponse
	for _, unspent := range unspents {
		rsp = append(rsp, proto.WalletUnspentResponse{
			Pos:          unspent.Pos,
			Txid:         unspent.Txid,
			Vout:         unspent.Vout,
			Value:        unspent.Value,
			Address:      unspent.Address,
			Confirmed:    unspent.Confirmed,
			SvrPubKey:    unspent.SvrPubKey,
			Scriptpubkey: unspent.Scriptpubkey,
		})
	}
	log.Info("api.wallet.unspent.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

func (h *Handler) walletTxs(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("walletTxs", r)
	if err != nil {
		log.Error("api.wallet.txs.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.WalletTxsRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].txs.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet.txs.req:%+v", req)

	if req.Limit > 256 {
		req.Limit = 256
	}
	txs, err := wdb.Txs(uid, req.Offset, req.Limit)
	if err != nil {
		log.Error("api.wallet[%v].txs.error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	var rsp []proto.WalletTxsResponse
	for _, tx := range txs {
		rsp = append(rsp, proto.WalletTxsResponse{
			Txid:        tx.Txid,
			Fee:         tx.Fee,
			Data:        tx.Data,
			Link:        tx.Link,
			Value:       tx.Value,
			Confirmed:   tx.Confirmed,
			BlockTime:   tx.BlockTime,
			BlockHeight: tx.BlockHeight,
		})
	}
	log.Info("api.wallet.txs.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

func (h *Handler) walletAddresses(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("walletAddresses", r)
	if err != nil {
		log.Error("api.wallet.addresses.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.WalletAddressesRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].addresses.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet[%v].addresses.req:%+v", uid, req)

	if req.Limit > 256 {
		req.Limit = 256
	}
	addresses, err := wdb.Addresses(uid, req.Offset, req.Limit)
	if err != nil {
		log.Error("api.wallet[%v].addresses.error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	var rsp []proto.WalletAddressesResponse
	for _, address := range addresses {
		rsp = append(rsp, proto.WalletAddressesResponse{
			Pos:     address.Pos,
			Address: address.Address,
		})
	}
	log.Info("api.wallet.addresses.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

func (h *Handler) walletSendFees(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("walletSendFees", r)
	if err != nil {
		log.Error("api.wallet.send.fees.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.WalletSendFeesRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].send.fees.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet[%v].send.fees.req:%+v", uid, req)

	fees, err := wdb.SendFees(uid, req.Priority, req.SendValue)
	if err != nil {
		log.Error("api.wallet[%v].send.fees.wdb.send.fees.error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	rsp := &proto.WalletSendFeesResponse{
		Fees:          fees.Fees,
		TotalValue:    fees.TotalValue,
		SendableValue: fees.SendableValue,
	}
	log.Info("api.wallet.fees.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

func (h *Handler) walletPortfolio(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)
	code := "CNY"

	// UID.
	uid, err := h.userinfo("walletPortfolio", r)
	if err != nil {
		log.Error("api.wallet.portfolio.uid.error:%+v", err)
		resp.writeError(err)
		return
	}

	// Request.
	req := &proto.WalletPortfolioRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Error("api.wallet[%v].portfolio.decode.body.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	log.Info("api.wallet[%v].portfolio.req:%+v", uid, req)

	if req.Code != "" {
		code = req.Code
	}
	ticker, err := wdb.store.getTicker(code)
	if err != nil {
		log.Error("api.wallet[%v].get.ticker.error:%+v", uid, err)
		resp.writeError(err)
		return
	}

	rsp := &proto.WalletPortfolioResponse{
		CoinSymbol:   "BTC",
		FiatSymbol:   ticker.Symbol,
		CurrentPrice: ticker.Last,
	}
	log.Info("api.wallet.portfolio.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}

func (h *Handler) walletPushTx(w http.ResponseWriter, r *http.Request) {
	log := h.log
	wdb := h.wdb
	resp := newResponse(log, w)

	// UID.
	uid, err := h.userinfo("walletPushTx", r)
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
	log.Info("api.wallet[%v].push.tx.req:%+v", uid, req)

	txid, err := wdb.chain.PushTx(req.TxHex)
	if err != nil {
		log.Error("api.wallet[%v].push.tx.error:%+v", uid, err)
		resp.writeError(err)
		return
	}
	rsp := &proto.TxPushResponse{
		TxID: txid,
	}
	log.Info("api.wallet.push.tx.rsp:%+v", rsp)
	resp.writeJSON(rsp)
}
