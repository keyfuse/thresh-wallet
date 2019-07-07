// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xcore"
	"github.com/tokublock/tokucore/xcore/bip32"
	"github.com/tokublock/tokucore/xcrypto"
)

func createSharedAddress(pos uint32, svrMasterPrvKey string, cliMasterPubkey string, net *network.Network) (string, error) {
	svrmasterkey, err := bip32.NewHDKeyFromString(svrMasterPrvKey)
	if err != nil {
		return "", err
	}
	svrchild, err := svrmasterkey.Derive(pos)
	if err != nil {
		return "", err
	}
	climasterkey, err := bip32.NewHDKeyFromString(cliMasterPubkey)
	if err != nil {
		return "", err
	}
	clichild, err := climasterkey.Derive(pos)
	if err != nil {
		return "", err
	}
	party := xcrypto.NewEcdsaParty(svrchild.PrivateKey())
	sharepub := party.Phase1(clichild.PublicKey())
	shared := xcore.NewPayToPubKeyHashAddress(sharepub.Hash160())
	return shared.ToString(net), nil
}
