// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package library

import (
	"testing"

	"server"

	"github.com/stretchr/testify/assert"
)

func TestWalletBalance(t *testing.T) {
	var token string

	ts := server.MockServer()
	defer ts.Close()

	mobile := "10086"
	// Token.
	{
		body := APIGetToken(ts.URL, mobile, "vcode", mockMasterPubKey)
		rsp := &TokenResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
		token = rsp.Token
	}

	body := APIWalletBalance(ts.URL, token)
	rsp := &WalletBalanceResponse{}
	unmarshal(body, rsp)

	t.Logf("%+v", rsp)
	assert.Equal(t, 200, rsp.Code)
}

func TestAPIEcdsaNewAddress(t *testing.T) {
	var token string

	ts := server.MockServer()
	defer ts.Close()

	mobile := "10086"
	// Token.
	{
		body := APIGetToken(ts.URL, mobile, "vcode", mockMasterPubKey)
		rsp := &TokenResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
		token = rsp.Token
	}

	body := APIEcdsaNewAddress(ts.URL, token)
	rsp := &EcdsaAddressResponse{}
	unmarshal(body, rsp)

	t.Logf("%+v", rsp)
	assert.Equal(t, 200, rsp.Code)
}
