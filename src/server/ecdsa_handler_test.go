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

func TestEcdsaAddressHandler(t *testing.T) {
	ts := MockServer()
	defer ts.Close()

	// Token.
	{
		req := &proto.TokenRequest{
			UID:          "bohu",
			MasterPubKey: "tpubD6NzVbkrYhZ4X7Cn1qGQ7XReumN4yFvgP3ms8dPtTiLD7wpP95cqmbaAkk5WSZaSrgpgtmPQhpNGmkxVRezP3WN486xEddsWHU22a6F7yJZ",
		}
		httpRsp, err := proto.NewRequest().Post(ts.URL+"/api/token", req)
		assert.Nil(t, err)

		resp := &proto.TokenResponse{}
		httpRsp.Json(resp)
		t.Log(resp)
	}

	// New address.
	{
		req := &proto.EcdsaAddressRequest{}
		httpRsp, err := proto.NewRequest().SetHeaders("Authorization", mockToken).Post(ts.URL+"/api/ecdsa/newaddress", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}
}
