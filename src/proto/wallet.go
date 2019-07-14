// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package proto

import ()

// WalletBalance --
type WalletBalanceResponse struct {
	AllBalance         uint64 `json:"all_balance"`
	UnconfirmedBalance uint64 `json:"confirmed_balance"`
}

// WalletUnspentRequest --
type WalletUnspentRequest struct {
	Amount uint64 `json:"amount"`
}

// WalletUnspentResponse --
type WalletUnspentResponse struct {
	Pos          uint32 `json:"pos"`
	Txid         string `json:"txid"`
	Vout         uint32 `json:"vout"`
	Value        uint64 `json:"value"`
	Address      string `json:"address"`
	Confirmed    bool   `json:"confirmed"`
	SvrPubKey    string `json:"svrpubkey"`
	Scriptpubkey string `json:"scriptpubkey"`
}

// TxPushRequest --
type TxPushRequest struct {
	TxHex string `json:"txhex"`
}

// TxPushResponse --
type TxPushResponse struct {
	TxID string `json:"txid"`
}
