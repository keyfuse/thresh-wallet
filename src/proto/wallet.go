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
