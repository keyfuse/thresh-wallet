// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package library

import (
	"net/http"

	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xcore/bip32"
)

// MasterPrvKeyResponse --
type MasterPrvKeyResponse struct {
	Status
	MasterPrvKey string `json:"masterprvkey"`
}

// NewMasterKey -- used to generate a new random master key.
func NewMasterPrvKey(chainnet string) string {
	rsp := &MasterPrvKeyResponse{}
	rsp.Code = http.StatusOK

	net := network.TestNet
	switch chainnet {
	case MainNet:
		net = network.MainNet
	}
	mk, err := bip32.NewHDKeyRand()
	if err != nil {
		rsp.Message = err.Error()
		rsp.Code = http.StatusInternalServerError
		return marshal(rsp)
	}
	rsp.MasterPrvKey = mk.ToString(net)
	return marshal(rsp)
}

// MasterPubKeyResponse --
type MasterPubKeyResponse struct {
	Status
	MasterPubKey string `json:"masterpubkey"`
}

// GetMasterPubKey -- get the master public key by priv key.
func GetMasterPubKey(chainnet string, masterPrvKey string) string {
	rsp := &MasterPubKeyResponse{}
	rsp.Code = http.StatusOK

	net := network.TestNet
	switch chainnet {
	case MainNet:
		net = network.MainNet
	}
	mk, err := bip32.NewHDKeyFromString(masterPrvKey)
	if err != nil {
		rsp.Message = err.Error()
		rsp.Code = http.StatusInternalServerError
		return marshal(rsp)
	}
	rsp.MasterPubKey = mk.HDPublicKey().ToString(net)
	return marshal(rsp)
}
