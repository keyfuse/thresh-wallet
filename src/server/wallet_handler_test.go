// thresh-wallet
//
// Copyright 2019 by KeyFuse
//
// GPLv3 License

package server

import (
	"testing"

	"proto"

	"github.com/stretchr/testify/assert"
)

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
