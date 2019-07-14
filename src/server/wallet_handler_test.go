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
	ts := MockServer()
	defer ts.Close()

	// Balance.
	{
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/balance", nil)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}
}

func TestWalletUnspent(t *testing.T) {
	ts := MockServer()
	defer ts.Close()

	{
		req := &proto.WalletUnspentRequest{
			Amount: 66,
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/unspent", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}
}

func TestWalletPushTx(t *testing.T) {
	ts := MockServer()
	defer ts.Close()

	{
		req := &proto.TxPushRequest{
			TxHex: "0100000002df5944657823338c5efa96eb0c89c01070bbc42a4593919682cb8a44df5c8c0f000000006b483045022100d4fadf37412aea7e562540a207107ac63713b2f799334551f93d10e4e0be452d02202d06f60b4409286cac3718845a7593d2bad0aff127de713a5ddee0e8be7901f5012103bc9fdab345b0c13b910b4f59cd6316eebdb75b1516cd67568ed2cbbb6387595bffffffff4a05631db2c351b26c21ccdfd09b2c4d2387da49b39ecee00799140db0b13523010000006b483045022100ec9235d80eb40cb52476fd9f78a016f2fa02d1c5a0f1fe27670fe3d2baf07a8502203e82a10877ea7bb9b22018c5f20c8baaa113e67c2b5ea4e2374f6de88bc2a921012103bc9fdab345b0c13b910b4f59cd6316eebdb75b1516cd67568ed2cbbb6387595bffffffff02a0860100000000001976a9143e1f197af8dd3cd441669ed34d096a0fbb06ec5388acda080000000000001976a914490e0eebcc5d462221ea38d00a6aee1238db2a5788ac00000000",
		}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/wallet/pushtx", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}
}
