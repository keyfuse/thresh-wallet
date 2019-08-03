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

func TestWalletCheck(t *testing.T) {
	var token string

	ts, cleanup := server.MockServer()
	defer cleanup()

	// Token.
	{
		body := APIGetToken(ts.URL, mockMobile1, "vcode")
		rsp := &TokenResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
		token = rsp.Token
	}

	body := APIWalletCheck(ts.URL, token)
	rsp := &WalletCheckResponse{}
	unmarshal(body, rsp)

	t.Logf("%+v", body)
	assert.Equal(t, 200, rsp.Code)
	assert.False(t, rsp.WalletExists)
	assert.False(t, rsp.BackupExists)
}

func TestWalletCreate(t *testing.T) {
	var token string

	ts, cleanup := server.MockServer()
	defer cleanup()

	// Token.
	{
		body := APIGetToken(ts.URL, mockMobile1, "vcode")
		rsp := &TokenResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
		token = rsp.Token
	}

	body := APIWalletCreate(ts.URL, token, mockMasterPrvKey)
	rsp := &WalletCreateResponse{}
	unmarshal(body, rsp)

	t.Logf("%+v", body)
	assert.Equal(t, 200, rsp.Code)
}

func TestWalletPortfolio(t *testing.T) {
	var token string

	ts, cleanup := server.MockServer()
	defer cleanup()

	// Token.
	{
		body := APIGetToken(ts.URL, mockMobile, "vcode")
		rsp := &TokenResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
		token = rsp.Token
	}

	body := APIWalletPortfolio(ts.URL, token)
	rsp := &WalletPortfolioResponse{}
	unmarshal(body, rsp)

	t.Logf("%+v", body)
	assert.Equal(t, 200, rsp.Code)
}

func TestWalletBalance(t *testing.T) {
	var token string

	ts, cleanup := server.MockServer()
	defer cleanup()

	// Token.
	{
		body := APIGetToken(ts.URL, mockMobile, "vcode")
		rsp := &TokenResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
		token = rsp.Token
	}

	body := APIWalletBalance(ts.URL, token)
	rsp := &WalletBalanceResponse{}
	unmarshal(body, rsp)

	t.Logf("%+v", body)
	assert.Equal(t, 200, rsp.Code)
	assert.Equal(t, uint64(103266), rsp.CoinValue)
}

func TestWalletTxs(t *testing.T) {
	var token string

	ts, cleanup := server.MockServer()
	defer cleanup()

	// Token.
	{
		body := APIGetToken(ts.URL, mockMobile, "vcode")
		rsp := &TokenResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
		token = rsp.Token
	}

	body := APIWalletTxs(ts.URL, token, 0, 2)
	rsp := &WalletTxsResponse{}
	unmarshal(body, rsp)

	t.Logf("%+v", body)
	assert.Equal(t, 200, rsp.Code)
	assert.Equal(t, 2, len(rsp.Txs))
}

func TestWalletAddresses(t *testing.T) {
	var token string

	ts, cleanup := server.MockServer()
	defer cleanup()

	// Token.
	{
		body := APIGetToken(ts.URL, mockMobile, "vcode")
		rsp := &TokenResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
		token = rsp.Token
	}

	body := APIWalletAddresses(ts.URL, token, 0, 2)
	rsp := &WalletAddressesResponse{}
	unmarshal(body, rsp)

	t.Logf("%+v", body)
	assert.Equal(t, 200, rsp.Code)
	assert.Equal(t, 2, len(rsp.Addresses))
}

func TestAPIEcdsaNewAddress(t *testing.T) {
	var token string

	ts, cleanup := server.MockServer()
	defer cleanup()

	// Token.
	{
		body := APIGetToken(ts.URL, mockMobile, "vcode")
		rsp := &TokenResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
		token = rsp.Token
	}

	for i := 0; i < 3; i++ {
		body := APIEcdsaNewAddress(ts.URL, token)
		rsp := &EcdsaAddressResponse{}
		unmarshal(body, rsp)

		t.Logf("%+v", body)
		assert.Equal(t, 200, rsp.Code)
	}
}

func TestAPIWalletSendFees(t *testing.T) {
	var token string

	ts, cleanup := server.MockServer()
	defer cleanup()

	// Token.
	{
		body := APIGetToken(ts.URL, mockMobile, "vcode")
		rsp := &TokenResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
		token = rsp.Token
	}

	{
		body := APIWalletSendFees(ts.URL, token, 100000)
		rsp := &WalletSendFeesResponse{}
		unmarshal(body, rsp)

		t.Logf("%+v", body)
		assert.Equal(t, 200, rsp.Code)
	}
}

func TestAPIWalletSend(t *testing.T) {
	var token string

	ts, cleanup := server.MockServer()
	defer cleanup()

	// Token.
	{
		body := APIGetToken(ts.URL, mockMobile, "vcode")
		rsp := &TokenResponse{}
		unmarshal(body, rsp)
		assert.Equal(t, 200, rsp.Code)
		token = rsp.Token
	}

	{
		body := APIWalletSend(ts.URL, token, "testnet", mockMasterPrvKey, "mmBRSnFG7o1BX5DaK8Da3xKxvjBh6fzNQq", 100000, 1000, "")
		rsp := &WalletSendResponse{}
		unmarshal(body, rsp)

		t.Logf("%+v", body)
		assert.Equal(t, 200, rsp.Code)
	}

	// Suffient value.
	{
		body := APIWalletSend(ts.URL, token, "testnet", mockMasterPrvKey, "mmBRSnFG7o1BX5DaK8Da3xKxvjBh6fzNQq", 1000000, 1000, "")
		rsp := &WalletSendResponse{}
		unmarshal(body, rsp)

		t.Logf("%+v", rsp)
		assert.Equal(t, 500, rsp.Code)
	}
}
