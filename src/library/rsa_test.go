// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package library

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRSAKey(t *testing.T) {
	var prvkey string
	var encrypted string
	msg := "Hello KeyFuse Labs"

	// Prv.
	{
		body := NewRSAPrvKey()
		t.Logf("body:%+v", body)

		rsp := &RSAKeyResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
		prvkey = rsp.PrvKey
	}

	// Pub.
	{
		body := GetRSAPubKey(prvkey)
		rsp := &RSAPubKeyResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
	}

	// Prv hash.
	{
		body := RSAPrvKeyHash(prvkey)
		t.Logf("body:%+v", body)

		rsp := &RSAPrvKeyHashResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
	}

	// Enc.
	{
		body := RSAEncrypt(msg, prvkey)
		t.Logf("body:%+v", body)

		rsp := &RSAEncryptResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
		encrypted = rsp.CipherText
	}

	// Dec.
	{
		body := RSADecrypt(encrypted, prvkey)
		t.Logf("body:%+v", body)

		rsp := &RSADecryptResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
	}

	// Sign.
	{
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte("sign text")))
		body := RSASign(hash, prvkey)
		t.Logf("body:%+v", body)

		rsp := &RSASignResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
	}
}
