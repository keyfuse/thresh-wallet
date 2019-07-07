// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package proto

import (
	"math/big"

	"github.com/tokublock/tokucore/xcrypto/paillier"
	"github.com/tokublock/tokucore/xcrypto/secp256k1"
)

// EcdsaAddressRequest --
type EcdsaAddressRequest struct {
	DeviceID string `json:"deviceid"`
}

// EcdsaAddressResponse --
type EcdsaAddressResponse struct {
	Pos     uint32 `json:"pos"`
	Address string `json:"address"`
}

// EcdsaR2Request --
type EcdsaR2Request struct {
	Pos  uint32            `json:"pos"`
	Hash []byte            `json:"hash"`
	R1   *secp256k1.Scalar `json:"R1"`
}

// EcdsaR2Response --
type EcdsaR2Response struct {
	Pos    uint32            `json:"pos"`
	R2     *secp256k1.Scalar `json:"R2"`
	ShareR *secp256k1.Scalar `json:"shareR"`
}

// EcdsaS2Request --
type EcdsaS2Request struct {
	Pos     uint32            `json:"pos"`
	Hash    []byte            `json:"hash"`
	EncPK1  *big.Int          `json:"encpk1"`
	EncPub1 *paillier.PubKey  `json:"encpub1"`
	R1      *secp256k1.Scalar `json:"R1"`
	ShareR  *secp256k1.Scalar `json:"shareR"`
}

// EcdsaS2Response --
type EcdsaS2Response struct {
	Pos uint32   `json:"pos"`
	S2  *big.Int `json:"S2"`
}
