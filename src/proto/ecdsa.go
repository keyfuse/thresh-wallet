// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package proto

import (
	"math/big"

	"github.com/keyfuse/tokucore/xcrypto/paillier"
	"github.com/keyfuse/tokucore/xcrypto/secp256k1"
)

// EcdsaR2Request --
type EcdsaR2Request struct {
	Pos  uint32            `json:"pos"`
	Hash []byte            `json:"hash"`
	R1   *secp256k1.Scalar `json:"R1"`
}

// EcdsaR2Response --
type EcdsaR2Response struct {
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
	S2 *big.Int `json:"S2"`
}
