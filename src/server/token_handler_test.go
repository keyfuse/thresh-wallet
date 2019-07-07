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

func TestVCodeHandler(t *testing.T) {
	ts := MockServer()
	defer ts.Close()

	// VCode.
	{
		req := &proto.VCodeRequest{
			UID: "bohu",
		}
		httpRsp, err := proto.NewRequest().Post(ts.URL+"/api/vcode", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}
}

func TestTokenHandler(t *testing.T) {
	ts := MockServer()
	defer ts.Close()

	// OK.
	{
		req := &proto.TokenRequest{
			UID:          "bohu",
			MasterPubKey: "tpubD6NzVbkrYhZ4X7Cn1qGQ7XReumN4yFvgP3ms8dPtTiLD7wpP95cqmbaAkk5WSZaSrgpgtmPQhpNGmkxVRezP3WN486xEddsWHU22a6F7yJZ",
		}

		httpRsp, err := proto.NewRequest().Post(ts.URL+"/api/token", req)
		assert.Nil(t, err)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}
}
