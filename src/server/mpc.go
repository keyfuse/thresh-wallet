// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package server

import (
	"fmt"
	"math/big"
	"strings"

	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"

	"github.com/keyfuse/tokucore/network"
	"github.com/keyfuse/tokucore/xcore"
	"github.com/keyfuse/tokucore/xcore/bip32"
	"github.com/keyfuse/tokucore/xcrypto"
	"github.com/keyfuse/tokucore/xcrypto/paillier"
	"github.com/keyfuse/tokucore/xcrypto/secp256k1"
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
		shared = xcore.NewPayToWitnessV0PubKeyHashAddress(sharepub.Hash160())
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

func rsaVerify(pubkeypem string, digestHex string, signatureHex string) error {
	digest, err := hex.DecodeString(digestHex)
	if err != nil {
		return err
	}

	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return err
	}

	block, _ := pem.Decode([]byte(pubkeypem))
	if block == nil {
		return fmt.Errorf("rsa.pubkey.pem.broken")
	}
	pubkey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(pubkey, crypto.SHA256, digest, signature)
}

func verifyPubKey(masterpubkey string, signatureHex string) error {
	hdpub, err := bip32.NewHDKeyFromString(masterpubkey)
	if err != nil {
		return err
	}
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return err
	}

	hash := sha256.Sum256([]byte(masterpubkey))
	return xcrypto.EcdsaVerify(hdpub.PublicKey(), hash[:], signature)
}
