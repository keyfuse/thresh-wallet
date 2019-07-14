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
			UID: mockUID,
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
			UID:          mockUID,
			MasterPubKey: mockCliMasterPubKey,
		}

		httpRsp, err := proto.NewRequest().Post(ts.URL+"/api/token", req)
		assert.Nil(t, err)
		rsp := &proto.TokenResponse{}
		httpRsp.Json(rsp)
		t.Logf("rsp:%+v", rsp)
		assert.Equal(t, 200, httpRsp.StatusCode())
	}
}
