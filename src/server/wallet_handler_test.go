// thresh-wallet
//
// Copyright 2019 by KeyFuse Labs
//
// GPLv3 License

package server

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"proto"

	"github.com/keyfuse/tokucore/network"
	"github.com/keyfuse/tokucore/xcore/bip32"
	"github.com/keyfuse/tokucore/xcrypto"

	"github.com/stretchr/testify/assert"
)

func TestNewAddressHandler(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	// Token.
	{
		req := &proto.TokenRequest{
			UID: mockUID,
		}
		httpRsp, err := proto.NewRequest().Post(ts.URL+"/api/token", req)
		assert.Nil(t, err)

		resp := &proto.TokenResponse{}
		httpRsp.Json(resp)
		t.Log(resp)
	}

	// New address.
	{
		req := &proto.WalletNewAddressRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/newaddress", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}
}

func TestWalletCheck(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	// Check.
	{
		req := &proto.WalletCheckRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/check", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		rsp := proto.WalletCheckResponse{}
		httpRsp.Json(&rsp)
		t.Logf("%v", rsp)
	}
}

func TestWalletCreate(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	var token string
	var signature string
	var masterpubkeywif string
	// Token.
	{
		req := &proto.TokenRequest{
			UID: mockEmail,
		}

		httpRsp, err := proto.NewRequest().Post(ts.URL+"/api/login/token", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		rsp := &proto.TokenResponse{}
		httpRsp.Json(rsp)
		token = rsp.Token
	}

	{
		hdkey, err := bip32.NewHDKeyRand()
		assert.Nil(t, err)
		pub := hdkey.HDPublicKey()
		masterpubkeywif = pub.ToString(network.TestNet)
		hash := sha256.Sum256([]byte(masterpubkeywif))
		sig, err := xcrypto.EcdsaSign(hdkey.PrivateKey(), hash[:])
		assert.Nil(t, err)
		signature = fmt.Sprintf("%x", sig)
	}

	// Create signature err.
	{
		req := &proto.WalletCreateRequest{
			Signature:    "16" + signature[2:],
			MasterPubKey: masterpubkeywif,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(ts.URL+"/api/wallet/create", req)
		assert.Nil(t, err)
		assert.Equal(t, 400, httpRsp.StatusCode())

		rsp := proto.WalletCreateResponse{}
		httpRsp.Json(&rsp)
		t.Logf("%v", rsp)
	}

	// Create ok.
	{
		req := &proto.WalletCreateRequest{
			Signature:    signature,
			MasterPubKey: masterpubkeywif,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(ts.URL+"/api/wallet/create", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		rsp := proto.WalletCreateResponse{}
		httpRsp.Json(&rsp)
		t.Logf("%v", rsp)
	}

	// Create exists error.
	{
		req := &proto.WalletCreateRequest{
			Signature:    signature,
			MasterPubKey: masterpubkeywif,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", token).Post(ts.URL+"/api/wallet/create", req)
		assert.Nil(t, err)
		assert.Equal(t, 500, httpRsp.StatusCode())
	}
}

func TestWalletBalance(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	// Balance.
	{
		req := &proto.WalletBalanceRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/balance", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}
}

func TestWalletUnspent(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	{
		req := &proto.WalletUnspentRequest{
			Amount: 66,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/unspent", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}
}

func TestWalletTxs(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	{
		req := &proto.WalletTxsRequest{
			Offset: 0,
			Limit:  3,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/txs", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		resp := []proto.WalletTxsResponse{}
		httpRsp.Json(&resp)
		assert.Equal(t, 3, len(resp))
	}

	// 2
	{
		req := &proto.WalletTxsRequest{
			Offset: 1,
			Limit:  3,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/txs", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		resp := []proto.WalletTxsResponse{}
		httpRsp.Json(&resp)
		assert.Equal(t, 2, len(resp))
	}

	// 1
	{
		req := &proto.WalletTxsRequest{
			Offset: 2,
			Limit:  3,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/txs", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		resp := []proto.WalletTxsResponse{}
		httpRsp.Json(&resp)
		assert.Equal(t, 1, len(resp))
	}

	// nil
	{
		req := &proto.WalletTxsRequest{
			Offset: 3,
			Limit:  3,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/txs", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		resp := []proto.WalletTxsResponse{}
		httpRsp.Json(&resp)
		assert.Equal(t, 0, len(resp))
	}
}

func TestWalletAddresses(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	{
		req := &proto.WalletAddressesRequest{
			Offset: 0,
			Limit:  3,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/addresses", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		resp := []proto.WalletTxsResponse{}
		httpRsp.Json(&resp)
		assert.Equal(t, 3, len(resp))
	}
}

func TestWalletSendFees(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	// Fast.
	{
		req := &proto.WalletSendFeesRequest{
			Priority:  "fast",
			SendValue: 1000,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/sendfees", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		got := &proto.WalletSendFeesResponse{}
		httpRsp.Json(got)

		want := &proto.WalletSendFeesResponse{
			Fees:          uint64(225),
			TotalValue:    uint64(103266),
			SendableValue: uint64(1000),
		}
		assert.Equal(t, want, got)
	}

	// Normal.
	{
		req := &proto.WalletSendFeesRequest{
			Priority:  "normal",
			SendValue: 1000,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/sendfees", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		got := &proto.WalletSendFeesResponse{}
		httpRsp.Json(got)

		want := &proto.WalletSendFeesResponse{
			Fees:          uint64(180),
			TotalValue:    uint64(103266),
			SendableValue: uint64(1000),
		}
		assert.Equal(t, want, got)
	}

	// Slow.
	{
		req := &proto.WalletSendFeesRequest{
			Priority:  "slow",
			SendValue: 1000,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/sendfees", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		got := &proto.WalletSendFeesResponse{}
		httpRsp.Json(got)

		want := &proto.WalletSendFeesResponse{
			Fees:          uint64(135),
			TotalValue:    uint64(103266),
			SendableValue: uint64(1000),
		}
		assert.Equal(t, want, got)
	}

	// Not enough.
	{
		req := &proto.WalletSendFeesRequest{
			Priority:  "fast",
			SendValue: 103260,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/sendfees", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		got := &proto.WalletSendFeesResponse{}
		httpRsp.Json(got)

		want := &proto.WalletSendFeesResponse{
			Fees:          uint64(374),
			TotalValue:    uint64(103266),
			SendableValue: uint64(102892),
		}
		assert.Equal(t, want, got)
	}

	// Send all.
	{
		req := &proto.WalletSendFeesRequest{
			Priority:  "fast",
			SendValue: 103266,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/sendfees", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		got := &proto.WalletSendFeesResponse{}
		httpRsp.Json(got)

		want := &proto.WalletSendFeesResponse{
			Fees:          uint64(374),
			TotalValue:    uint64(103266),
			SendableValue: uint64(102892),
		}
		assert.Equal(t, want, got)
	}
}

func TestWalletPortfolio(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	// Code empty.
	{
		req := &proto.WalletPortfolioRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/portfolio", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		got := &proto.WalletPortfolioResponse{}
		httpRsp.Json(got)

		want := &proto.WalletPortfolioResponse{
			CoinSymbol:   "BTC",
			FiatSymbol:   "¥",
			CurrentPrice: 73711.13,
		}
		assert.Equal(t, want, got)
	}

	// Code USD.
	{
		req := &proto.WalletPortfolioRequest{
			Code: "USD",
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/portfolio", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())

		got := &proto.WalletPortfolioResponse{}
		httpRsp.Json(got)

		want := &proto.WalletPortfolioResponse{
			CoinSymbol:   "BTC",
			FiatSymbol:   "$",
			CurrentPrice: 10721.13,
		}
		assert.Equal(t, want, got)
	}
}

func TestWalletPushTx(t *testing.T) {
	ts, cleanup := MockServer()
	defer cleanup()

	{
		req := &proto.TxPushRequest{
			TxHex: "0100000002df5944657823338c5efa96eb0c89c01070bbc42a4593919682cb8a44df5c8c0f000000006b483045022100d4fadf37412aea7e562540a207107ac63713b2f799334551f93d10e4e0be452d02202d06f60b4409286cac3718845a7593d2bad0aff127de713a5ddee0e8be7901f5012103bc9fdab345b0c13b910b4f59cd6316eebdb75b1516cd67568ed2cbbb6387595bffffffff4a05631db2c351b26c21ccdfd09b2c4d2387da49b39ecee00799140db0b13523010000006b483045022100ec9235d80eb40cb52476fd9f78a016f2fa02d1c5a0f1fe27670fe3d2baf07a8502203e82a10877ea7bb9b22018c5f20c8baaa113e67c2b5ea4e2374f6de88bc2a921012103bc9fdab345b0c13b910b4f59cd6316eebdb75b1516cd67568ed2cbbb6387595bffffffff02a0860100000000001976a9143e1f197af8dd3cd441669ed34d096a0fbb06ec5388acda080000000000001976a914490e0eebcc5d462221ea38d00a6aee1238db2a5788ac00000000",
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/pushtx", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}
}
