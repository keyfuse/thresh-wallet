// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package library

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"net/http"
)

const (
	bits = 2048
)

func strToRSAPrvKey(prvkey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(prvkey))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// RSAKeyResponse --
type RSAKeyResponse struct {
	Status
	PrvKey string `json:"prvkey"`
}

// NewRSAPrvKey -- generates a new key pair.
func NewRSAPrvKey() string {
	rsp := &RSAKeyResponse{}
	rsp.Code = http.StatusOK

	prvkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		rsp.Message = err.Error()
		rsp.Code = http.StatusInternalServerError
		return marshal(rsp)
	}
	rsp.PrvKey = string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(prvkey)}))
	return marshal(rsp)
}

// RSAPubKeyResponse --
type RSAPubKeyResponse struct {
	Status
	PubKey string `json:"pubkey"`
}

// GetRSAPubKey -- used to get the pem format of the pubkey.
func GetRSAPubKey(rsaPrvKey string) string {
	rsp := &RSAPubKeyResponse{}
	rsp.Code = http.StatusOK

	prv, err := strToRSAPrvKey(rsaPrvKey)
	if err != nil {
		rsp.Message = err.Error()
		rsp.Code = http.StatusInternalServerError
		return marshal(rsp)
	}
	rsp.PubKey = string(pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&prv.PublicKey)}))
	return marshal(rsp)
}

// RSAPrvKeyHashResponse --
type RSAPrvKeyHashResponse struct {
	Status
	Hash string `json:"hash"`
}

// RSAPrvKeyHash -- used to get the sha256 of the pubkey with pem.
func RSAPrvKeyHash(prvkey string) string {
	rsp := &RSAPrvKeyHashResponse{}
	rsp.Code = http.StatusOK

	rsp.Hash = fmt.Sprintf("%x", (sha256.Sum256([]byte(prvkey))))
	return marshal(rsp)
}

// RSAEncryptResponse --
type RSAEncryptResponse struct {
	Status
	CipherText string `json:"ciphertext"`
}

// RSAEncrypt -- encrypt the msg with prvkey.
func RSAEncrypt(msg string, prvkey string) string {
	rsp := &RSAEncryptResponse{}
	rsp.Code = http.StatusOK

	prv, err := strToRSAPrvKey(prvkey)
	if err != nil {
		rsp.Message = err.Error()
		rsp.Code = http.StatusInternalServerError
		return marshal(rsp)
	}
	enc, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &prv.PublicKey, []byte(msg), nil)
	if err != nil {
		rsp.Message = err.Error()
		rsp.Code = http.StatusInternalServerError
		return marshal(rsp)
	}
	rsp.CipherText = hex.EncodeToString(enc)
	return marshal(rsp)
}

// RSADecryptResponse --
type RSADecryptResponse struct {
	Status
	PlainText string `json:"plaintext"`
}

// RSADecrypt -- decrypt the cipher to the plain with prvkey.
func RSADecrypt(ciphertext string, prvkey string) string {
	rsp := &RSADecryptResponse{}
	rsp.Code = http.StatusOK

	cipher, err := hex.DecodeString(ciphertext)
	if err != nil {
		rsp.Message = err.Error()
		rsp.Code = http.StatusInternalServerError
		return marshal(rsp)
	}

	prv, err := strToRSAPrvKey(prvkey)
	if err != nil {
		rsp.Message = err.Error()
		rsp.Code = http.StatusInternalServerError
		return marshal(rsp)
	}

	text, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, prv, cipher, nil)
	if err != nil {
		rsp.Message = err.Error()
		rsp.Code = http.StatusInternalServerError
		return marshal(rsp)
	}
	rsp.PlainText = string(text)
	return marshal(rsp)
}

// RSASignResponse --
type RSASignResponse struct {
	Status
	Signature string `json:"signature"`
}

// RSASign -- sign the digest with prvkey.
func RSASign(digestHex string, prvkey string) string {
	rsp := &RSASignResponse{}
	rsp.Code = http.StatusOK

	digest, err := hex.DecodeString(digestHex)
	if err != nil {
		rsp.Message = err.Error()
		rsp.Code = http.StatusInternalServerError
		return marshal(rsp)
	}

	prv, err := strToRSAPrvKey(prvkey)
	if err != nil {
		rsp.Message = err.Error()
		rsp.Code = http.StatusInternalServerError
		return marshal(rsp)
	}

	sig, err := rsa.SignPKCS1v15(rand.Reader, prv, crypto.SHA256, digest)
	if err != nil {
		rsp.Message = err.Error()
		rsp.Code = http.StatusInternalServerError
		return marshal(rsp)
	}
	rsp.Signature = fmt.Sprintf("%x", sig)
	return marshal(rsp)
}
