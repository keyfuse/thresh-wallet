// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xcore"
	"github.com/tokublock/tokucore/xcore/bip32"
	"github.com/tokublock/tokucore/xcrypto"
	"github.com/tokublock/tokucore/xcrypto/paillier"
	"github.com/tokublock/tokucore/xcrypto/secp256k1"
)

const ()

func createSharedAddress(pos uint32, svrMasterPrvKey string, cliMasterPubkey string, net *network.Network, typ string) (string, error) {
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

	var shared xcore.Address
	switch strings.ToUpper(typ) {
	case "P2PKH":
		shared = xcore.NewPayToPubKeyHashAddress(sharepub.Hash160())
	default:
		shared = xcore.NewPayToWitnessPubKeyHashAddress(sharepub.Hash160())
	}
	return shared.ToString(net), nil
}

// createEcdsaR2 -- used to create the R2.
// Returns:
// R2, ShareR
func createEcdsaR2(pos uint32, svrMasterPrvKey string, hash []byte, R1 *secp256k1.Scalar) (*secp256k1.Scalar, *secp256k1.Scalar, error) {
	masterkey, err := bip32.NewHDKeyFromString(svrMasterPrvKey)
	if err != nil {
		return nil, nil, err
	}
	childkey, err := masterkey.Derive(pos)
	if err != nil {
		return nil, nil, err
	}

	bobParty := xcrypto.NewEcdsaParty(childkey.PrivateKey())
	// Skip phase1.
	// Phase2.
	_, _, r2 := bobParty.Phase2(hash)
	shareR := bobParty.Phase3(R1)
	return r2, shareR, nil
}

// createEcdsaS2 -- used to create S2.
// Returns:
// S2
func createEcdsaS2(pos uint32, svrMasterPrvKey string, hash []byte, R1 *secp256k1.Scalar, shareR *secp256k1.Scalar, encPK1 *big.Int, encPub1 *paillier.PubKey) (*big.Int, error) {
	masterkey, err := bip32.NewHDKeyFromString(svrMasterPrvKey)
	if err != nil {
		return nil, err
	}
	childkey, err := masterkey.Derive(pos)
	if err != nil {
		return nil, err
	}

	bobParty := xcrypto.NewEcdsaParty(childkey.PrivateKey())
	// Skip phase1.
	// Phase2.
	bobParty.Phase2(hash)
	bobShareR := bobParty.Phase3(R1)
	if bobShareR.X.Cmp(shareR.X) != 0 || bobShareR.Y.Cmp(shareR.Y) != 0 {
		return nil, fmt.Errorf("api.ecdsa.s2.shareR.not.equal")
	}
	// Skip phase3.
	// Phase4.
	return bobParty.Phase4(encPK1, encPub1, shareR)
}

func createSvrChildPubKey(pos uint32, svrMasterPrvKey string, net *network.Network) (string, error) {
	svrmasterkey, err := bip32.NewHDKeyFromString(svrMasterPrvKey)
	if err != nil {
		return "", err
	}
	svrchild, err := svrmasterkey.Derive(pos)
	if err != nil {
		return "", err
	}
	return svrchild.HDPublicKey().ToString(net), nil
}
