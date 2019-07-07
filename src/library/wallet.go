// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package library

import (
	"fmt"
	"net/http"

	"proto"
)

// WalletBalanceResponse --
type WalletBalanceResponse struct {
	Status
	AllBalance         uint64 `json:"all_balance"`
	UnconfirmedBalance uint64 `json:"confirmed_balance"`
}

// APIWalletBalance -- Wallet balance api.
func APIWalletBalance(url string, token string) string {
	rsp := &WalletBalanceResponse{}
	rsp.Code = http.StatusOK
	path := fmt.Sprintf("%s/api/wallet/balance", url)

	httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, nil)
	if err != nil {
		rsp.Code = http.StatusInternalServerError
		rsp.Message = err.Error()
		return marshal(rsp)
	}

	balance := &proto.WalletBalanceResponse{}
	if err := httpRsp.Json(balance); err != nil {
		rsp.Code = httpRsp.StatusCode()
		rsp.Message = err.Error()
		return marshal(rsp)
	}
	rsp.AllBalance = balance.AllBalance
	rsp.UnconfirmedBalance = balance.UnconfirmedBalance
	return marshal(rsp)
}

// EcdsaAddressResponse --
type EcdsaAddressResponse struct {
	Status
	Pos     uint32 `json:"pos"`
	Address string `json:"address"`
}

// APIEcdsaNewAddress -- ecdsa new address api.
func APIEcdsaNewAddress(url string, token string) string {
	rsp := &EcdsaAddressResponse{}
	rsp.Code = http.StatusOK
	path := fmt.Sprintf("%s/api/ecdsa/newaddress", url)

	req := &proto.EcdsaAddressRequest{}

	httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(path, req)
	if err != nil {
		rsp.Code = http.StatusInternalServerError
		rsp.Message = err.Error()
		return marshal(rsp)
	}

	address := &proto.EcdsaAddressResponse{}
	if err := httpRsp.Json(address); err != nil {
		rsp.Code = http.StatusInternalServerError
		rsp.Message = err.Error()
		return marshal(rsp)
	}
	rsp.Pos = address.Pos
	rsp.Address = address.Address
	return marshal(rsp)
}
